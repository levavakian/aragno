package component

import (
	"reflect"
)

// SpiderBody tags an entity as a spider body
type SpiderBody struct {
	Name string
}

var SpiderBodyType = reflect.TypeOf(&SpiderBody{})

// SpiderLeg tags an entity as a spider leg
type SpiderLeg struct {
}

var SpiderLegType = reflect.TypeOf(&SpiderLeg{})
