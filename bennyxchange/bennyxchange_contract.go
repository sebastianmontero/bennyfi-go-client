package bennyxchange

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	eos "github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
	"github.com/sebastianmontero/eos-go-toolbox/util"
	"github.com/sebastianmontero/eos-go/ecc"
)

var (
	ExchangeEntry             = eos.Name("xentry")
	ExchangeBeneficiaryReward = eos.Name("xbenefreward")
	ExchangePoolManagerFee    = eos.Name("xpoolmgrfee")

	OfferTypeBuy  = eos.Name("buy")
	OfferTypeSell = eos.Name("sell")

	SettingBatchSize = "BATCH_SIZE"
)

type BaseOffer struct {
	Seller           eos.AccountName `json:"seller"`
	Buyer            eos.AccountName `json:"buyer"`
	OfferType        eos.Name        `json:"offer_type"`
	ItemID           uint64          `json:"item_id"`
	Price            eos.Asset       `json:"price"`
	PoolName         string          `json:"pool_name"`
	PoolType         eos.Name        `json:"pool_type"`
	PoolStakeEndTime eos.TimePoint   `json:"pool_stake_end_time"`
}

func (m *BaseOffer) IsBuyOffer() bool {
	return m.OfferType == OfferTypeBuy
}

func (m *BaseOffer) IsSellOffer() bool {
	return m.OfferType == OfferTypeSell
}

type Offer struct {
	OfferID uint64 `json:"offer_id"`
	*BaseOffer
	ExpirationTime eos.TimePoint `json:"expiration_time"`
	CreatedDate    eos.TimePoint `json:"created_date"`
}

func (m *Offer) ToMakeOfferArgs(exchangeType eos.Name) *MakeOfferArgs {
	var who eos.AccountName
	if m.OfferType == OfferTypeBuy {
		who = m.Buyer
	} else {
		who = m.Seller
	}
	return &MakeOfferArgs{
		Who:            who,
		OfferType:      m.OfferType,
		ExchangeType:   exchangeType,
		ItemID:         m.ItemID,
		Price:          m.Price,
		ExpirationTime: m.ExpirationTime,
	}
}

func (m *Offer) ToAcceptedOffer() *AcceptedOffer {
	return &AcceptedOffer{
		BaseOffer: m.BaseOffer,
	}
}

func (m *Offer) Lapse() {
	m.ExpirationTime = util.ShiftedTimePoint(-1 * time.Millisecond)
}

func (m *Offer) String() string {
	result, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed marshalling offer: %v", err))
	}
	return string(result)
}

type AcceptedOffer struct {
	AcceptedOfferID uint64 `json:"accepted_offer_id"`
	*BaseOffer
	AcceptedDate eos.TimePoint `json:"accepted_date"`
}

func (m *AcceptedOffer) String() string {
	result, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed marshalling accepted offer: %v", err))
	}
	return string(result)
}

type MakeOfferArgs struct {
	Who            eos.AccountName `json:"who"`
	OfferType      eos.Name        `json:"offer_type"`
	ExchangeType   eos.Name        `json:"exchange_type"`
	ItemID         uint64          `json:"item_id"`
	Price          eos.Asset       `json:"price"`
	ExpirationTime eos.TimePoint   `json:"expiration_time"`
}

type BennyXchangeContract struct {
	*contract.SettingsContract
	callCounter uint64
}

func NewBennyXchangeContract(eos *service.EOS, contractName string) *BennyXchangeContract {
	return &BennyXchangeContract{
		contract.NewSettingsContract(eos, contractName),
		0,
	}
}

func (m *BennyXchangeContract) NextCallCounter() uint64 {
	m.callCounter++
	return m.callCounter
}

func (m *BennyXchangeContract) ExecAction(permissionLevel interface{}, action string, actionData interface{}) (string, error) {
	resp, err := m.Contract.ExecAction(permissionLevel, action, actionData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Tx ID: %v", resp.TransactionID), nil
}

func (m *BennyXchangeContract) ConfigureOpenPermission(publicKey *ecc.PublicKey) error {
	openActions := []string{
		"delexpirdoff",
	}
	err := m.EOS.CreateSimplePermission(m.ContractName, "open", publicKey)
	if err != nil {
		return fmt.Errorf("failed to create open permission, error: %v", err)
	}
	for _, action := range openActions {
		err = m.EOS.LinkPermission(m.ContractName, action, "open", false)
		if err != nil {
			return fmt.Errorf("failed to link open permission to the %v action, error: %v", action, err)
		}
	}
	return nil
}

func (m *BennyXchangeContract) MakeOffer(offer *Offer, exchangeType eos.Name) (string, error) {
	return m.MakeOfferFromArgs(offer.ToMakeOfferArgs(exchangeType))
}

func (m *BennyXchangeContract) MakeOfferWithAuthorizer(offer *Offer, exchangeType eos.Name, authorizer eos.AccountName) (string, error) {
	return m.MakeOfferFromArgsWithAuthorizer(offer.ToMakeOfferArgs(exchangeType), authorizer)
}

func (m *BennyXchangeContract) MakeOfferFromArgs(args *MakeOfferArgs) (string, error) {
	return m.ExecAction(args.Who, "makeoffer", args)
}

func (m *BennyXchangeContract) MakeOfferFromArgsWithAuthorizer(args *MakeOfferArgs, authorizer eos.AccountName) (string, error) {
	return m.ExecAction(authorizer, "makeoffer", args)
}

func (m *BennyXchangeContract) AcceptOffer(who eos.AccountName, exchangeType eos.Name, offerId uint64) (string, error) {
	actionData := struct {
		Who          eos.AccountName
		ExchangeType eos.Name
		OfferId      uint64
	}{who, exchangeType, offerId}
	return m.ExecAction(who, "acceptoffer", actionData)
}

func (m *BennyXchangeContract) AcceptOfferWithAuthorizer(authorizer, who eos.AccountName, exchangeType eos.Name, offerId uint64) (string, error) {
	actionData := struct {
		Who          eos.AccountName
		ExchangeType eos.Name
		OfferId      uint64
	}{who, exchangeType, offerId}
	return m.ExecAction(authorizer, "acceptoffer", actionData)
}

func (m *BennyXchangeContract) DeleteOffer(who eos.AccountName, exchangeType eos.Name, offerId uint64) (string, error) {
	actionData := struct {
		Who          eos.AccountName
		ExchangeType eos.Name
		OfferId      uint64
	}{who, exchangeType, offerId}
	return m.ExecAction(who, "deleteoffer", actionData)
}

func (m *BennyXchangeContract) DeleteOfferWithAuthorizer(authorizer, who eos.AccountName, exchangeType eos.Name, offerId uint64) (string, error) {
	actionData := struct {
		Who          eos.AccountName
		ExchangeType eos.Name
		OfferId      uint64
	}{who, exchangeType, offerId}
	return m.ExecAction(authorizer, "deleteoffer", actionData)
}

func (m *BennyXchangeContract) DeleteExpiredOffers(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "delexpirdoff", callCounter)
}

func (m *BennyXchangeContract) Reset(limit uint64, toDelete []string) (string, error) {
	actionData := struct {
		Limit       uint64
		ToDelete    []string
		CallCounter uint64
	}{limit, toDelete, m.NextCallCounter()}
	return m.ExecAction(eos.AN(m.ContractName), "reset", actionData)
}

func (m *BennyXchangeContract) TstLapseTime(exchangeType eos.Name, offerId uint64) (string, error) {
	actionData := struct {
		ExchangeType eos.Name
		OfferId      uint64
		CallCounter  uint64
	}{exchangeType, offerId, m.NextCallCounter()}
	return m.ExecAction(eos.AN(m.ContractName), "tstlapsetime", actionData)
}

func (m *BennyXchangeContract) GetOffer(exchangeType eos.Name, offerId uint64) (*Offer, error) {
	offers, err := m.GetOffersReq(&eos.GetTableRowsRequest{
		Scope:      exchangeType.String(),
		LowerBound: strconv.FormatUint(offerId, 10),
		UpperBound: strconv.FormatUint(offerId, 10),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(offers) > 0 {
		return &offers[0], nil
	}
	return nil, nil
}

func (m *BennyXchangeContract) GetAllOffersAsMap() ([]map[string]interface{}, error) {
	return m.GetAllTableRowsWithScopesAsMap("offers", "offer_id", "0", nil)
}

func (m *BennyXchangeContract) GetAllAcceptedOffersAsMap() ([]map[string]interface{}, error) {
	return m.GetAllTableRowsWithScopesAsMap("acceptdoffrs", "accepted_offer_id", "0", nil)
}

func (m *BennyXchangeContract) GetAcceptedOffer(exchangeType eos.Name, acceptedOfferId uint64) (*AcceptedOffer, error) {
	acceptedOffers, err := m.GetAcceptedOffersReq(&eos.GetTableRowsRequest{
		Scope:      exchangeType.String(),
		LowerBound: strconv.FormatUint(acceptedOfferId, 10),
		UpperBound: strconv.FormatUint(acceptedOfferId, 10),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(acceptedOffers) > 0 {
		return &acceptedOffers[0], nil
	}
	return nil, nil
}

func (m *BennyXchangeContract) GetLastOffer(exchangeType eos.Name) (*Offer, error) {
	offers, err := m.GetOffersReq(&eos.GetTableRowsRequest{
		Scope:   exchangeType.String(),
		Reverse: true,
		Limit:   1,
	})
	if err != nil {
		return nil, err
	}
	if len(offers) > 0 {
		return &offers[0], nil
	}
	return nil, nil
}

func (m *BennyXchangeContract) GetLastAcceptedOffer(exchangeType eos.Name) (*AcceptedOffer, error) {
	acceptedOffers, err := m.GetAcceptedOffersReq(&eos.GetTableRowsRequest{
		Scope:   exchangeType.String(),
		Reverse: true,
		Limit:   1,
	})
	if err != nil {
		return nil, err
	}
	if len(acceptedOffers) > 0 {
		return &acceptedOffers[0], nil
	}
	return nil, nil
}

func (m *BennyXchangeContract) GetOffersReq(req *eos.GetTableRowsRequest) ([]Offer, error) {

	var offers []Offer
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "offers"
	err := m.GetTableRows(*req, &offers)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return offers, nil
}

func (m *BennyXchangeContract) GetAcceptedOffersReq(req *eos.GetTableRowsRequest) ([]AcceptedOffer, error) {

	var acceptedOffers []AcceptedOffer
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "acceptdoffrs"
	err := m.GetTableRows(*req, &acceptedOffers)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return acceptedOffers, nil
}

func (m *BennyXchangeContract) GetOffersBySellerAndId(seller eos.AccountName) ([]Offer, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterOffersBySellerAndId(request, seller)
	if err != nil {
		return nil, err
	}
	return m.GetOffersReq(request)
}

func (m *BennyXchangeContract) FilterOffersBySellerAndId(req *eos.GetTableRowsRequest, seller eos.AccountName) error {

	req.Index = "3"
	req.KeyType = "i128"
	req.Reverse = true
	lb, err := m.EOS.GetComposedIndexValue(seller, 0)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	ub, err := m.EOS.GetComposedIndexValue(seller, uint64(18446744073709551615))
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	// fmt.Println("LB: ", lb, "UB: ", ub)
	req.LowerBound = lb
	req.UpperBound = ub
	return err
}

func (m *BennyXchangeContract) GetOffersByBuyerAndId(buyer eos.AccountName) ([]Offer, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterOffersByBuyerAndId(request, buyer)
	if err != nil {
		return nil, err
	}
	return m.GetOffersReq(request)
}

func (m *BennyXchangeContract) FilterOffersByBuyerAndId(req *eos.GetTableRowsRequest, buyer eos.AccountName) error {

	req.Index = "4"
	req.KeyType = "i128"
	req.Reverse = true
	lb, err := m.EOS.GetComposedIndexValue(buyer, 0)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	ub, err := m.EOS.GetComposedIndexValue(buyer, uint64(18446744073709551615))
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	// fmt.Println("LB: ", lb, "UB: ", ub)
	req.LowerBound = lb
	req.UpperBound = ub
	return err
}

func (m *BennyXchangeContract) GetOffersByOfferTypeAndId(offerType eos.Name) ([]Offer, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterOffersByOfferTypeAndId(request, offerType)
	if err != nil {
		return nil, err
	}
	return m.GetOffersReq(request)
}

func (m *BennyXchangeContract) FilterOffersByOfferTypeAndId(req *eos.GetTableRowsRequest, offerType eos.Name) error {

	req.Index = "5"
	req.KeyType = "i128"
	req.Reverse = true
	lb, err := m.EOS.GetComposedIndexValue(offerType, 0)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	ub, err := m.EOS.GetComposedIndexValue(offerType, uint64(18446744073709551615))
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	// fmt.Println("LB: ", lb, "UB: ", ub)
	req.LowerBound = lb
	req.UpperBound = ub
	return err
}

func (m *BennyXchangeContract) GetOffersByItemIdAndSymbol(itemID uint64) ([]Offer, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterOffersByItemIdAndSymbol(request, itemID)
	if err != nil {
		return nil, err
	}
	return m.GetOffersReq(request)
}

func (m *BennyXchangeContract) FilterOffersByItemIdAndSymbol(req *eos.GetTableRowsRequest, itemID uint64) error {

	req.Index = "6"
	req.KeyType = "i128"
	req.Reverse = true
	lb, err := m.EOS.GetComposedIndexValue(itemID, 0)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	ub, err := m.EOS.GetComposedIndexValue(itemID, uint64(18446744073709551615))
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	// fmt.Println("LB: ", lb, "UB: ", ub)
	req.LowerBound = lb
	req.UpperBound = ub
	return err
}

func (m *BennyXchangeContract) GetAcceptedOffersBySellerAndId(seller eos.AccountName) ([]AcceptedOffer, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterAcceptedOffersBySellerAndId(request, seller)
	if err != nil {
		return nil, err
	}
	return m.GetAcceptedOffersReq(request)
}

func (m *BennyXchangeContract) FilterAcceptedOffersBySellerAndId(req *eos.GetTableRowsRequest, seller eos.AccountName) error {

	req.Index = "3"
	req.KeyType = "i128"
	req.Reverse = true
	lb, err := m.EOS.GetComposedIndexValue(seller, 0)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	ub, err := m.EOS.GetComposedIndexValue(seller, uint64(18446744073709551615))
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	// fmt.Println("LB: ", lb, "UB: ", ub)
	req.LowerBound = lb
	req.UpperBound = ub
	return err
}

func (m *BennyXchangeContract) GetAcceptedOffersByBuyerAndId(buyer eos.AccountName) ([]AcceptedOffer, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterAcceptedOffersByBuyerAndId(request, buyer)
	if err != nil {
		return nil, err
	}
	return m.GetAcceptedOffersReq(request)
}

func (m *BennyXchangeContract) FilterAcceptedOffersByBuyerAndId(req *eos.GetTableRowsRequest, buyer eos.AccountName) error {

	req.Index = "4"
	req.KeyType = "i128"
	req.Reverse = true
	lb, err := m.EOS.GetComposedIndexValue(buyer, 0)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	ub, err := m.EOS.GetComposedIndexValue(buyer, uint64(18446744073709551615))
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	// fmt.Println("LB: ", lb, "UB: ", ub)
	req.LowerBound = lb
	req.UpperBound = ub
	return err
}

func (m *BennyXchangeContract) GetAcceptedOffersByOfferTypeAndId(offerType eos.Name) ([]AcceptedOffer, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterAcceptedOffersByOfferTypeAndId(request, offerType)
	if err != nil {
		return nil, err
	}
	return m.GetAcceptedOffersReq(request)
}

func (m *BennyXchangeContract) FilterAcceptedOffersByOfferTypeAndId(req *eos.GetTableRowsRequest, offerType eos.Name) error {

	req.Index = "5"
	req.KeyType = "i128"
	req.Reverse = true
	lb, err := m.EOS.GetComposedIndexValue(offerType, 0)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	ub, err := m.EOS.GetComposedIndexValue(offerType, uint64(18446744073709551615))
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	// fmt.Println("LB: ", lb, "UB: ", ub)
	req.LowerBound = lb
	req.UpperBound = ub
	return err
}

func (m *BennyXchangeContract) GetAcceptedOffersByItemIdAndId(itemID uint64) ([]AcceptedOffer, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterAcceptedOffersByItemIdAndSymbol(request, itemID)
	if err != nil {
		return nil, err
	}
	return m.GetAcceptedOffersReq(request)
}

func (m *BennyXchangeContract) FilterAcceptedOffersByItemIdAndSymbol(req *eos.GetTableRowsRequest, itemID uint64) error {

	req.Index = "6"
	req.KeyType = "i128"
	req.Reverse = true
	lb, err := m.EOS.GetComposedIndexValue(itemID, 0)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	ub, err := m.EOS.GetComposedIndexValue(itemID, uint64(18446744073709551615))
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	// fmt.Println("LB: ", lb, "UB: ", ub)
	req.LowerBound = lb
	req.UpperBound = ub
	return err
}
