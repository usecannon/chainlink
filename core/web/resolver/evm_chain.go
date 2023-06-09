package resolver

import (
	"context"

	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/web/loader"
)

// ChainResolver resolves the Chain type.
type ChainResolver struct {
	chain types.ChainConfig
}

func NewChain(chain types.ChainConfig) *ChainResolver {
	return &ChainResolver{chain: chain}
}

func NewChains(chains []types.ChainConfig) []*ChainResolver {
	var resolvers []*ChainResolver
	for _, c := range chains {
		resolvers = append(resolvers, NewChain(c))
	}

	return resolvers
}

// ID resolves the chain's unique identifier.
func (r *ChainResolver) ID() graphql.ID {
	return graphql.ID(r.chain.ID.String())
}

// Enabled resolves the chain's enabled field.
func (r *ChainResolver) Enabled() bool {
	return r.chain.Enabled
}

// Config resolves the chain's configuration field
func (r *ChainResolver) Config() *ChainConfigResolver {
	return NewChainConfig(*r.chain.Cfg)
}

func (r *ChainResolver) Nodes(ctx context.Context) ([]*NodeResolver, error) {
	nodes, err := loader.GetNodesByChainID(ctx, r.chain.ID.String())
	if err != nil {
		return nil, err
	}

	return NewNodes(nodes), nil
}

type ChainPayloadResolver struct {
	chain types.ChainConfig
	NotFoundErrorUnionType
}

func NewChainPayload(chain types.ChainConfig, err error) *ChainPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "chain not found", isExpectedErrorFn: nil}

	return &ChainPayloadResolver{chain: chain, NotFoundErrorUnionType: e}
}

func (r *ChainPayloadResolver) ToChain() (*ChainResolver, bool) {
	if r.err != nil {
		return nil, false
	}

	return NewChain(r.chain), true
}

type ChainsPayloadResolver struct {
	chains []types.ChainConfig
	total  int32
}

func NewChainsPayload(chains []types.ChainConfig, total int32) *ChainsPayloadResolver {
	return &ChainsPayloadResolver{chains: chains, total: total}
}

func (r *ChainsPayloadResolver) Results() []*ChainResolver {
	return NewChains(r.chains)
}

func (r *ChainsPayloadResolver) Metadata() *PaginationMetadataResolver {
	return NewPaginationMetadata(r.total)
}

// -- CreateChain Mutation --

type CreateChainPayloadResolver struct {
	chain     *types.ChainConfig
	inputErrs map[string]string
}

func NewCreateChainPayload(chain *types.ChainConfig, inputErrs map[string]string) *CreateChainPayloadResolver {
	return &CreateChainPayloadResolver{chain: chain, inputErrs: inputErrs}
}

func (r *CreateChainPayloadResolver) ToCreateChainSuccess() (*CreateChainSuccessResolver, bool) {
	if r.chain == nil {
		return nil, false
	}

	return NewCreateChainSuccess(r.chain), true
}

func (r *CreateChainPayloadResolver) ToInputErrors() (*InputErrorsResolver, bool) {
	if r.inputErrs != nil {
		var errs []*InputErrorResolver

		for path, message := range r.inputErrs {
			errs = append(errs, NewInputError(path, message))
		}

		return NewInputErrors(errs), true
	}

	return nil, false
}

type CreateChainSuccessResolver struct {
	chain *types.ChainConfig
}

func NewCreateChainSuccess(chain *types.ChainConfig) *CreateChainSuccessResolver {
	return &CreateChainSuccessResolver{chain: chain}
}

func (r *CreateChainSuccessResolver) Chain() *ChainResolver {
	return NewChain(*r.chain)
}

type UpdateChainPayloadResolver struct {
	chain     *types.ChainConfig
	inputErrs map[string]string
	NotFoundErrorUnionType
}

func NewUpdateChainPayload(chain *types.ChainConfig, inputErrs map[string]string, err error) *UpdateChainPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "chain not found", isExpectedErrorFn: nil}

	return &UpdateChainPayloadResolver{chain: chain, inputErrs: inputErrs, NotFoundErrorUnionType: e}
}

func (r *UpdateChainPayloadResolver) ToUpdateChainSuccess() (*UpdateChainSuccessResolver, bool) {
	if r.chain == nil {
		return nil, false
	}

	return NewUpdateChainSuccess(*r.chain), true
}

func (r *UpdateChainPayloadResolver) ToInputErrors() (*InputErrorsResolver, bool) {
	if r.inputErrs != nil {
		var errs []*InputErrorResolver

		for path, message := range r.inputErrs {
			errs = append(errs, NewInputError(path, message))
		}

		return NewInputErrors(errs), true
	}

	return nil, false
}

type UpdateChainSuccessResolver struct {
	chain types.ChainConfig
}

func NewUpdateChainSuccess(chain types.ChainConfig) *UpdateChainSuccessResolver {
	return &UpdateChainSuccessResolver{chain: chain}
}

func (r *UpdateChainSuccessResolver) Chain() *ChainResolver {
	return NewChain(r.chain)
}

type DeleteChainPayloadResolver struct {
	chain *types.ChainConfig
	NotFoundErrorUnionType
}

func NewDeleteChainPayload(chain *types.ChainConfig, err error) *DeleteChainPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "chain not found", isExpectedErrorFn: nil}

	return &DeleteChainPayloadResolver{chain: chain, NotFoundErrorUnionType: e}
}

func (r *DeleteChainPayloadResolver) ToDeleteChainSuccess() (*DeleteChainSuccessResolver, bool) {
	if r.chain == nil {
		return nil, false
	}

	return NewDeleteChainSuccess(*r.chain), true
}

type DeleteChainSuccessResolver struct {
	chain types.ChainConfig
}

func NewDeleteChainSuccess(chain types.ChainConfig) *DeleteChainSuccessResolver {
	return &DeleteChainSuccessResolver{chain: chain}
}

func (r *DeleteChainSuccessResolver) Chain() *ChainResolver {
	return NewChain(r.chain)
}
