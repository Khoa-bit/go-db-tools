//go:generate go run tools/gen.go
package main

import (
	"fmt"
	"go-db-tools/tool"
	"log"
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

// -------------------------------------------------------------------

var _ NestedModeler = (*Layer0)(nil)

type Layer0 struct {
	ID       int64
	Layers1  []Layer1
	Layers12 []Layer2
}

func (l *Layer0) GetID() int64 {
	return l.ID
}

func (l *Layer0) Append(modeler Modeler) {
	switch model := modeler.(type) {
	case *Layer1:
		l.Layers1 = append(l.Layers1, *model)
	case *Layer2:
		l.Layers12 = append(l.Layers12, *model)
	default:
		tool.Assert(false, "failed to append modeler to Layer0", "modeler", modeler)
	}
}

// -------------------------------------------------------------------

var _ NestedModeler = (*Layer1)(nil)

type Layer1 struct {
	ID      int64
	Layers2 []Layer2
}

func (l *Layer1) GetID() int64 {
	return l.ID
}

func (l *Layer1) Append(modeler Modeler) {
	layer2, ok := modeler.(*Layer2)
	tool.Assert(ok, "failed to append modeler to Layer1", "modeler", modeler)
	l.Layers2 = append(l.Layers2, *layer2)
}

// -------------------------------------------------------------------

var _ Modeler = (*Layer2)(nil)

type Layer2 struct {
	ID int64
}

func (l *Layer2) GetID() int64 {
	return l.ID
}

// -------------------------------------------------------------------

var _ NestedModeler = (*emptyModel)(nil)

type emptyModel struct {
}

func (e *emptyModel) GetID() int64 {
	return 0
}

func (e *emptyModel) Append(modeler Modeler) {
	log.Fatalf("cannot append to emptyModel: '%v'", modeler)
}

// -------------------------------------------------------------------

type Row struct {
	Layer0  Layer0
	Layer1  Layer1
	Layer12 Layer2
	Layer2  Layer2
}

func main() {
	rows := []Row{
		{Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 201}, Layer12: Layer2{ID: 5301}, Layer2: Layer2{ID: 301}},
		{Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 201}, Layer12: Layer2{ID: 5302}, Layer2: Layer2{ID: 302}},
		{Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 201}, Layer12: Layer2{ID: 5303}, Layer2: Layer2{ID: 303}},
		{Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 204}, Layer12: Layer2{ID: 5304}, Layer2: Layer2{ID: 304}},
		{Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 205}, Layer12: Layer2{ID: 5305}, Layer2: Layer2{ID: 305}},
		{Layer0: Layer0{ID: 106}, Layer1: Layer1{ID: 206}, Layer12: Layer2{ID: 5306}, Layer2: Layer2{ID: 306}},
		{Layer0: Layer0{ID: 106}, Layer1: Layer1{ID: 206}, Layer12: Layer2{ID: 5307}, Layer2: Layer2{ID: 307}},
		{Layer0: Layer0{ID: 106}, Layer1: Layer1{ID: 206}, Layer12: Layer2{ID: 5308}, Layer2: Layer2{ID: 308}},
		{Layer0: Layer0{ID: 109}, Layer1: Layer1{ID: 209}, Layer12: Layer2{ID: 5000}, Layer2: Layer2{ID: 000}},
		{Layer0: Layer0{ID: 110}, Layer1: Layer1{ID: 000}, Layer12: Layer2{ID: 5000}, Layer2: Layer2{ID: 000}},
	}

	builder := NestedModelBuilder{}
	builder12 := NestedModelBuilder{}
	for _, row := range rows {
		builder.Build([]NestedModeler{&row.Layer0, &row.Layer1}, &row.Layer2)
		builder12.Build([]NestedModeler{&row.Layer0}, &row.Layer12)
	}
	results := GetAll[*Layer0](builder, builder12)
	log.Printf("Results: \n%s\n", tool.DebugMarshal(results))
	log.Println("=====================================")
	rows2 := []Row{
		{Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 201}, Layer12: Layer2{ID: 5301}, Layer2: Layer2{ID: 301}},
		{Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 201}, Layer12: Layer2{ID: 5302}, Layer2: Layer2{ID: 302}},
		{Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 201}, Layer12: Layer2{ID: 5303}, Layer2: Layer2{ID: 303}},
		{Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 204}, Layer12: Layer2{ID: 5304}, Layer2: Layer2{ID: 304}},
		{Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 205}, Layer12: Layer2{ID: 5305}, Layer2: Layer2{ID: 305}},
		{Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 206}, Layer12: Layer2{ID: 5306}, Layer2: Layer2{ID: 306}},
		{Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 206}, Layer12: Layer2{ID: 5307}, Layer2: Layer2{ID: 307}},
		{Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 206}, Layer12: Layer2{ID: 5308}, Layer2: Layer2{ID: 308}},
		{Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 209}, Layer12: Layer2{ID: 5000}, Layer2: Layer2{ID: 000}},
		{Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 000}, Layer12: Layer2{ID: 5000}, Layer2: Layer2{ID: 000}},
	}

	builder2 := NestedModelBuilder{}
	builder22 := NestedModelBuilder{}
	for _, row := range rows2 {
		builder2.Build([]NestedModeler{&row.Layer0, &row.Layer1}, &row.Layer2)
		builder22.Build([]NestedModeler{&row.Layer0}, &row.Layer12)
	}
	result, ok := GetOne[*Layer0](builder2, builder22)
	log.Printf("Result: \n%s\nok: %v\n", tool.DebugMarshal(result), ok)
}
