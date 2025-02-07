package patch

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"terraform-provider-identitynow/internal/util"
	"unsafe"

	sailpointBeta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
)

const (
	jsonTag = "json"
)

type comparableValues struct {
	modifiedVal interface{}
	currentVal  interface{}
	path        string
}

type patchBuilder interface {
	GenerateJsonPatch() ([]sailpointBeta.JsonPatchOperation, error)
	defineValuesToCompare()
}

type abstractPatchBuilder struct {
	operations            []sailpointBeta.JsonPatchOperation
	valuesToCompare       []comparableValues
	referencesToCompare   []comparableValues
	defineValuesToCompare func()
}

func (pb *abstractPatchBuilder) GenerateJsonPatch() ([]sailpointBeta.JsonPatchOperation, error) {
	pb.defineValuesToCompare()
	if err := pb.doCompare(); err != nil {
		return nil, err
	}
	return pb.operations, nil
}

func (pb *abstractPatchBuilder) doCompare() error {
	if err := pb.compareValues(); err != nil {
		return err
	}
	if err := pb.compareReferences(); err != nil {
		return err
	}
	return nil
}
func (pb *abstractPatchBuilder) compareValues() error {
	if pb.valuesToCompare == nil || len(pb.valuesToCompare) == 0 {
		return nil
	}
	for _, value := range pb.valuesToCompare {
		if err := pb.compareValue(value.modifiedVal, value.currentVal, value.path); err != nil {
			return err
		}
	}
	return nil
}

func (pb *abstractPatchBuilder) compareReferences() error {
	if pb.referencesToCompare == nil || len(pb.referencesToCompare) == 0 {
		return nil
	}
	for _, value := range pb.referencesToCompare {
		if err := pb.compareInnerObject(value.modifiedVal, value.currentVal, value.path, "Id"); err != nil {
			return err
		}
	}
	return nil
}

func (pb *abstractPatchBuilder) compareValue(modifiedVal, currentVal interface{}, path string) error {
	modified := reflect.ValueOf(modifiedVal)
	current := reflect.ValueOf(currentVal)
	// the data structures must be identical
	if modifiedVal != nil && currentVal != nil && modified.Kind() != current.Kind() {
		return fmt.Errorf("kind does not match at: %s modified: %s current: %s", path, modified.Kind(), current.Kind())
	}
	switch modified.Kind() {
	case reflect.Struct:
		return pb.compareStruct(modifiedVal, currentVal, path)
	case reflect.Ptr:
		return pb.comparePointer(modified, current, path)
	case reflect.Slice, reflect.Array:
		return pb.compareArray(modifiedVal, currentVal, path)
	case reflect.Map:
		return pb.compareMap(modifiedVal.(map[string]interface{}), currentVal.(map[string]interface{}), path)
		//case reflect.Interface:
	//	return pb.processInterface(modified, current, pointer)
	case reflect.String:
		var modifiedString, currentString *string
		if modifiedVal != nil {
			val := modified.String()
			modifiedString = &val
		}
		if currentVal != nil {
			val := current.String()
			currentString = &val
		}
		return pb.compareString(modifiedString, currentString, path)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return pb.compareInt64(modified.Int(), current.Int(), path)
	case reflect.Bool:
		return pb.compareBool(modified.Bool(), current.Bool(), path)
	case reflect.Float32, reflect.Float64:
		return pb.compareFloat64(modified.Float(), current.Float(), path)
	case reflect.Invalid:
		if modifiedVal == nil && currentVal != nil {
			return pb.remove(path)
		}
		// undefined interfaces are ignored for now
		return nil
	default:
		return fmt.Errorf("compareValue unsupported kind: %s at: %s", modified.Kind(), path)
	}
}

func (pb *abstractPatchBuilder) comparePointer(modified reflect.Value, current reflect.Value, path string) error {
	if !modified.IsNil() && !current.IsNil() {
		// the values of the pointers will be processed in a next step
		if err := pb.compareValue(modified.Elem().Interface(), current.Elem().Interface(), path); err != nil {
			return err
		}
	} else if !modified.IsNil() {
		return pb.add(modified.Elem().Interface(), path)
	} else if !current.IsNil() {
		return pb.remove(path)
	}
	return nil
}

func (pb *abstractPatchBuilder) compareString(modified, current *string, path string) error {
	if modified == current || (modified != nil && current != nil && *modified == *current) {
		return nil
	}
	if modified == nil || *modified == "" {
		return pb.remove(path)
	} else if current == nil || *current == "" {
		return pb.add(*modified, path)
	} else {
		return pb.replace(*modified, path)
	}
}

func (pb *abstractPatchBuilder) compareStringPointer(modified, current *string, path string) error {
	if (modified == nil && current == nil) || (*modified == *current) {
		return nil
	}
	if modified == nil || *modified == "" {
		return pb.remove(path)
	} else if current == nil || *current == "" {
		return pb.add(*modified, path)
	} else {
		return pb.replace(*modified, path)
	}
}

func (pb *abstractPatchBuilder) compareInt64(modified, current int64, path string) error {
	if modified == current {
		return nil
	}
	if modified == 0 {
		return pb.remove(path)
	} else if current == 0 {
		return pb.add(modified, path)
	} else {
		return pb.replace(modified, path)
	}
}

func (pb *abstractPatchBuilder) compareInt64Pointer(modified, current *int64, path string) error {
	if (modified == nil && current == nil) || (*modified == *current) {
		return nil
	}
	if *modified == 0 {
		return pb.remove(path)
	} else if *current == 0 {
		return pb.add(*modified, path)
	} else {
		return pb.replace(*modified, path)
	}
}

func (pb *abstractPatchBuilder) compareInt32(modified, current int32, path string) error {
	if modified == current {
		return nil
	}
	if modified == 0 {
		return pb.remove(path)
	} else if current == 0 {
		return pb.add(modified, path)
	} else {
		return pb.replace(modified, path)
	}
}

func (pb *abstractPatchBuilder) compareFloat64(modified, current float64, path string) error {
	if modified == current {
		return nil
	}
	if modified == 0 {
		return pb.remove(path)
	} else if current == 0 {
		return pb.add(modified, path)
	} else {
		return pb.replace(modified, path)
	}
}

func (pb *abstractPatchBuilder) compareInt32Pointer(modified, current *int32, path string) error {
	if (modified == nil && current == nil) || (*modified == *current) {
		return nil
	}
	if *modified == 0 {
		return pb.remove(path)
	} else if *current == 0 {
		return pb.add(*modified, path)
	} else {
		return pb.replace(*modified, path)
	}
}

func (pb *abstractPatchBuilder) compareBool(modified, current bool, path string) error {
	if modified == current {
		return nil
	}
	return pb.replace(modified, path)
}

func (pb *abstractPatchBuilder) compareBoolPointer(modified, current *bool, path string) error {
	if *modified == *current {
		return nil
	}
	return pb.replace(*modified, path)
}

func (pb *abstractPatchBuilder) compareArray(modVal, curVal interface{}, path string) error {
	if modVal == nil {
		return pb.remove(path)
	}
	if curVal == nil {
		return pb.add(modVal, path)
	}
	if !reflect.DeepEqual(modVal, curVal) {
		return pb.replace(modVal, path)
	}
	return nil
}

func (pb *abstractPatchBuilder) compareMap(modAttr, curAttr map[string]interface{}, path string) error {
	if util.IsNil(modAttr) && util.IsNil(curAttr) {
		return nil
	}
	if !util.IsNil(modAttr) && util.IsNil(curAttr) {
		return pb.add(modAttr, path)
	}
	if util.IsNil(modAttr) && !util.IsNil(curAttr) {
		return pb.remove(path)
	}

	basePath := path + "/"
	for key, curValue := range curAttr {
		modValue, ok := modAttr[key]
		if !ok {
			if err := pb.remove(basePath + key); err != nil {
				return err
			}
		} else {
			if err := pb.compareValue(modValue, curValue, basePath+key); err != nil {
				return err
			}
		}
	}
	for key, modValue := range modAttr {
		_, ok := curAttr[key]
		if !ok {
			if err := pb.add(modValue, basePath+key); err != nil {
				return err
			}
		}
	}
	return nil
}

func (pb *abstractPatchBuilder) compareInnerObject(mod, cur interface{}, path, fieldName string) error {
	modified := pb.getInnerObject(mod)
	current := pb.getInnerObject(cur)
	if util.IsNil(modified) && util.IsNil(current) {
		return nil
	}
	if !util.IsNil(modified) && util.IsNil(current) {
		return pb.add(modified, path)
	}
	if util.IsNil(modified) && !util.IsNil(current) {
		return pb.remove(path)
	}
	idMod, err := pb.getFieldValue(modified, fieldName)
	if err != nil {
		return err
	}
	idCur, err := pb.getFieldValue(current, fieldName)
	if err != nil {
		return err
	}
	if idMod == idCur {
		return nil
	}
	return pb.replace(idMod, path+"/"+strings.ToLower(fieldName))
}

func (pb *abstractPatchBuilder) getInnerObject(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	ref := reflect.ValueOf(value)
	if ref.Kind() == reflect.Struct {
		ref2 := reflect.New(ref.Type()).Elem()
		ref2.Set(ref)
		isSet := ref2.FieldByName("isSet")
		innerObject := ref2.FieldByName("value")
		if innerObject.IsValid() && isSet.IsValid() {
			if innerObject.IsNil() {
				return nil
			}
			return GetUnexportedField(innerObject)
		}
		return value
	} else if ref.Kind() == reflect.Pointer {
		elem := ref.Elem()
		if !elem.IsValid() {
			return nil
		}
		return pb.getInnerObject(elem.Interface())
	}
	return value
}

func GetUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

func (pb *abstractPatchBuilder) getFieldValue(obj interface{}, fieldName string) (string, error) {
	ref := reflect.ValueOf(obj)
	elem := ref
	if ref.Kind() == reflect.Ptr {
		elem = ref.Elem()
	}
	id := elem.FieldByName(fieldName)
	switch id.Kind() {
	case reflect.Pointer:
		return id.Elem().String(), nil
	case reflect.String:
		return id.String(), nil
	default:
		return "", fmt.Errorf("getFieldValue unsupported type %v", id.Kind())
	}
}

func (pb *abstractPatchBuilder) remove(path string) error {
	pb.operations = append(pb.operations, sailpointBeta.JsonPatchOperation{
		Op:   "remove",
		Path: path,
	})
	return nil
}

func (pb *abstractPatchBuilder) add(value any, path string) error {
	operationValue, err := pb.valueToOperationValue(value, path)
	if err != nil {
		return err
	}
	pb.operations = append(pb.operations, sailpointBeta.JsonPatchOperation{
		Op:    "add",
		Path:  path,
		Value: &operationValue,
	})
	return nil
}

func (pb *abstractPatchBuilder) replace(value any, path string) error {
	operationValue, err := pb.valueToOperationValue(value, path)
	if err != nil {
		return err
	}
	pb.operations = append(pb.operations, sailpointBeta.JsonPatchOperation{
		Op:    "replace",
		Path:  path,
		Value: &operationValue,
	})
	return nil
}

func (pb *abstractPatchBuilder) valueToOperationValue(value any, path string) (sailpointBeta.UpdateMultiHostSourcesRequestInnerValue, error) {
	var opValue sailpointBeta.UpdateMultiHostSourcesRequestInnerValue
	ref := reflect.ValueOf(value)
	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}
	switch ref.Kind() {
	case reflect.String:
		v := value.(string)
		opValue = sailpointBeta.StringAsUpdateMultiHostSourcesRequestInnerValue(&v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := int32(ref.Int())
		opValue = sailpointBeta.Int32AsUpdateMultiHostSourcesRequestInnerValue(&v)
	case reflect.Float32, reflect.Float64:
		// Patch operation does not support Float, so converting to String
		v := strconv.FormatFloat(ref.Float(), 'f', -1, 64)
		opValue = sailpointBeta.StringAsUpdateMultiHostSourcesRequestInnerValue(&v)
	case reflect.Bool:
		v := value.(bool)
		opValue = sailpointBeta.BoolAsUpdateMultiHostSourcesRequestInnerValue(&v)
	case reflect.Map:
		v := value.(map[string]interface{})
		opValue = sailpointBeta.MapmapOfStringAnyAsUpdateMultiHostSourcesRequestInnerValue(&v)
	case reflect.Array, reflect.Slice:
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{}, err
		}
		var inner []sailpointBeta.ArrayInner
		err = json.Unmarshal(jsonBytes, &inner)
		if err != nil {
			return sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{}, err
		}
		opValue = sailpointBeta.ArrayOfArrayInnerAsUpdateMultiHostSourcesRequestInnerValue(&inner)
	case reflect.Struct:
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{}, err
		}
		var inner map[string]interface{}
		err = json.Unmarshal(jsonBytes, &inner)
		if err != nil {
			return sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{}, err
		}
		opValue = sailpointBeta.MapmapOfStringAnyAsUpdateMultiHostSourcesRequestInnerValue(&inner)
	default:
		return sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{}, fmt.Errorf("valueToOperationValue unsupported kind: %s at: %s", ref.Kind(), path)
	}
	return opValue, nil
}

func (pb *abstractPatchBuilder) compareStruct(modified, current interface{}, path string) error {
	if reflect.DeepEqual(modified, current) {
		return nil
	}

	if util.IsNil(modified) && !util.IsNil(current) {
		return pb.remove(path)
	}
	if !util.IsNil(modified) && util.IsNil(current) {
		return pb.add(modified, path)
	}
	modRef := reflect.ValueOf(modified)
	curRef := reflect.ValueOf(current)
	if pb.isNullabeStruct(modRef) {
		modValue := modRef.MethodByName("Get").Call([]reflect.Value{})[0]
		curValue := curRef.MethodByName("Get").Call([]reflect.Value{})[0]
		return pb.compareValue(modValue.Interface(), curValue.Interface(), path)
	} else {
		for j := 0; j < modRef.NumField(); j++ {
			field := modRef.Type().Field(j)
			tag := strings.Split(field.Tag.Get(jsonTag), ",")[0]
			if tag == "" || tag == "_" || !modRef.Field(j).CanInterface() {
				// struct fields without a JSON tag set or unexported fields are ignored
				continue
			}
			if err := pb.compareValue(modRef.Field(j).Interface(), curRef.Field(j).Interface(), path+"/"+tag); err != nil {
				return err
			}
		}
	}
	return nil
}

func (pb *abstractPatchBuilder) isNullabeStruct(structField reflect.Value) bool {
	if strings.HasPrefix(structField.Type().Name(), "Nullable") && structField.MethodByName("Get").IsValid() {
		for j := 0; j < structField.NumField(); j++ {
			field := structField.Type().Field(j)
			if field.Name == "value" {
				return true
			}
		}
	}
	return false
}
