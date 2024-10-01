package mapper

import (
	"fmt"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ToPayload(payload map[string]*anypb.Any) (map[string]any, error) {
	output := make(map[string]any)
	for key, anyValue := range payload {
		var value any
		var err error

		switch anyValue.TypeUrl {
		case "type.googleapis.com/google.protobuf.StringValue":
			var stringValue wrapperspb.StringValue
			if err = anyValue.UnmarshalTo(&stringValue); err != nil {
				return nil, fmt.Errorf("failed to unmarshal StringValue: %w", err)
			}
			value = stringValue.Value

		case "type.googleapis.com/google.protobuf.Int32Value":
			var intValue wrapperspb.Int32Value
			if err = anyValue.UnmarshalTo(&intValue); err != nil {
				return nil, fmt.Errorf("failed to unmarshal Int32Value: %w", err)
			}
			value = intValue.Value

		case "type.googleapis.com/google.protobuf.Int64Value":
			var intValue wrapperspb.Int64Value
			if err = anyValue.UnmarshalTo(&intValue); err != nil {
				return nil, fmt.Errorf("failed to unmarshal Int64Value: %w", err)
			}
			value = intValue.Value

		case "type.googleapis.com/google.protobuf.UInt32Value":
			var uintValue wrapperspb.UInt32Value
			if err = anyValue.UnmarshalTo(&uintValue); err != nil {
				return nil, fmt.Errorf("failed to unmarshal UInt32Value: %w", err)
			}
			value = uintValue.Value

		case "type.googleapis.com/google.protobuf.UInt64Value":
			var uintValue wrapperspb.UInt64Value
			if err = anyValue.UnmarshalTo(&uintValue); err != nil {
				return nil, fmt.Errorf("failed to unmarshal UInt64Value: %w", err)
			}
			value = uintValue.Value

		case "type.googleapis.com/google.protobuf.FloatValue":
			var floatValue wrapperspb.FloatValue
			if err = anyValue.UnmarshalTo(&floatValue); err != nil {
				return nil, fmt.Errorf("failed to unmarshal FloatValue: %w", err)
			}
			value = floatValue.Value

		case "type.googleapis.com/google.protobuf.DoubleValue":
			var doubleValue wrapperspb.DoubleValue
			if err = anyValue.UnmarshalTo(&doubleValue); err != nil {
				return nil, fmt.Errorf("failed to unmarshal DoubleValue: %w", err)
			}
			value = doubleValue.Value

		case "type.googleapis.com/google.protobuf.BoolValue":
			var boolValue wrapperspb.BoolValue
			if err = anyValue.UnmarshalTo(&boolValue); err != nil {
				return nil, fmt.Errorf("failed to unmarshal BoolValue: %w", err)
			}
			value = boolValue.Value

		case "type.googleapis.com/google.protobuf.BytesValue":
			var bytesValue wrapperspb.BytesValue
			if err = anyValue.UnmarshalTo(&bytesValue); err != nil {
				return nil, fmt.Errorf("failed to unmarshal BytesValue: %w", err)
			}
			value = bytesValue.Value

		default:
			return nil, fmt.Errorf("unsupported type URL: %s", anyValue.TypeUrl)
		}

		output[key] = value
	}
	return output, nil
}
