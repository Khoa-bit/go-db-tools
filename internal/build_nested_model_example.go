package internal

import (
	"go-db-tools/tool"
	"log"
)

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

func BuildNestedModelExample() {
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
