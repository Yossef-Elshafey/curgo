package eval

import "curgo/types/object"

func stringInterface(left *object.String, member string) object.Object {
	switch member {
	case "length":
		return &object.Integer{Value: int64(len(left.Value))}
	default: 
		return newError("object %s donest support operation %s", left.Type(), member)
	}
}

func responseInterface(left *object.Response, member string) object.Object {
	switch member {
	case "status":
		return &object.Integer{Value: int64(left.Res.StatusCode)}
	case "statusText":
		return &object.String{Value: left.Res.Status}
	case "cookies":
		// TODO:
		return newError("object %s doesn't support operation %s", left.Type(), member)
	case "location":
		// TODO:
		return newError("object %s doesn't support operation %s", left.Type(), member)
	case "header":
		// TODO:
		return newError("object %s doesn't support operation %s", left.Type(), member)
	case "body":
		// TODO:
		return newError("object %s doesn't support operation %s", left.Type(), member)
	default: 
		return newError("object %s doesn't support operation %s", left.Type(), member)
	}
}
