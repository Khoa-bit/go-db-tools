package internal

import (
	"fmt"
	"go-db-tools/tool"
	"time"
)

type Modeler interface {
	GetID() int64
}

var _ Modeler = (*Model)(nil)

type Model struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	DeletedBy string
}

func (m *Model) GetID() int64 {
	return m.ID
}

type NestedModeler interface {
	GetID() int64
	Append(modeler Modeler)
}

type NestedModelBuilder struct {
	Results   []NestedModeler
	Layers    []NestedModeler
	LayerLast Modeler
}

func GetAll[T NestedModeler](builders ...NestedModelBuilder) []T {
	tool.Assert(len(builders) > 0, "expected at least one builder")
	for i, b := range builders {
		if b.Results == nil || b.Layers == nil || b.LayerLast == nil {
			result := make([]T, 0)
			return result
		}
		if b.LayerLast.GetID() > 0 {
			b.Layers[len(b.Layers)-1].Append(b.LayerLast)
		}
		for l := len(b.Layers) - 1; l >= 1; l-- {
			if b.Layers[l].GetID() > 0 {
				b.Layers[l-1].Append(b.Layers[l])
			}
		}
		if b.Layers[0].GetID() > 0 {
			builders[i].Results = append(b.Results, b.Layers[0])
		}
	}
	results := make([]T, len(builders[0].Results))
	for i := 0; i < len(builders[0].Results); i++ {
		result, ok := any(builders[0].Results[i]).(T)
		tool.Assert(ok, fmt.Sprintf("expected type %T, got %T", result, builders[0].Results))
		results[i] = result
	}
	return results
}

func GetOne[T NestedModeler](builders ...NestedModelBuilder) (T, bool) {
	tool.Assert(len(builders) > 0, "expected at least one builder")
	for i, b := range builders {
		if b.Results == nil || b.Layers == nil || b.LayerLast == nil {
			var result T
			return result, false
		}
		if b.LayerLast.GetID() > 0 {
			b.Layers[len(b.Layers)-1].Append(b.LayerLast)
		}
		for l := len(b.Layers) - 1; l >= 1; l-- {
			if b.Layers[l].GetID() > 0 {
				b.Layers[l-1].Append(b.Layers[l])
			}
		}
		if b.Layers[0].GetID() > 0 {
			builders[i].Results = append(b.Results, b.Layers[0])
		}
	}
	count := len(builders[0].Results)
	if count == 0 {
		var result T
		return result, false
	}
	tool.Assert(count == 1, "expected at most one result", "count", count)
	result, ok := any(builders[0].Results[0]).(T)
	tool.Assert(ok, fmt.Sprintf("expected type %T, got %T", result, builders[0].Results))
	return result, true
}

func (b *NestedModelBuilder) Build(layers []NestedModeler, layerLast Modeler) {
	// Init builder
	if b.Results == nil || b.Layers == nil || b.LayerLast == nil {
		b.Results = make([]NestedModeler, 0, 10)
		b.Layers = make([]NestedModeler, len(layers))
		for i := 0; i < len(layers); i++ {
			b.Layers[i] = &emptyModel{}
		}
		b.LayerLast = &emptyModel{}
	}

	for i, layer := range layers {
		tool.Assert(layer != nil, "layer must not be nil", "i", i)
	}
	tool.Assert(layerLast != nil, "layerLast must not be nil")
	prevLayers := make([]NestedModeler, len(layers))
	for i, layer := range b.Layers {
		prevLayers[i] = layer
	}
	prevLayerLast := b.LayerLast
	flush := make([]bool, len(layers))
	flushLayerLast := false

	for i := 0; i < len(layers); i++ {
		forceFlush := false
		if i-1 >= 0 {
			forceFlush = flush[i-1]
		}

		if layers[i].GetID() > 0 && (forceFlush || prevLayers[i].GetID() != layers[i].GetID()) {
			// Store the current layer
			b.Layers[i] = layers[i]
			// Reset all next layers
			for j := i + 1; j < len(layers); j++ {
				b.Layers[j] = &emptyModel{}
			}
			b.LayerLast = &emptyModel{}
			// Flush all layers from the current layer
			for k := i; k < len(layers); k++ {
				flush[k] = true
			}
			flushLayerLast = true
		}
	}
	if layerLast.GetID() > 0 && (flushLayerLast || prevLayerLast.GetID() != layerLast.GetID()) {
		b.LayerLast = layerLast
		flushLayerLast = true
	}

	if flushLayerLast && prevLayerLast.GetID() > 0 {
		prevLayers[len(layers)-1].Append(prevLayerLast)
	}
	for i := len(layers) - 1; i >= 1; i-- {
		if flush[i] && prevLayers[i].GetID() > 0 {
			prevLayers[i-1].Append(prevLayers[i])
		}
	}
	if flush[0] && prevLayers[0].GetID() > 0 {
		b.Results = append(b.Results, prevLayers[0])
	}
}
