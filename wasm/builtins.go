package wasm

import (
	"github.com/deckarep/golang-set"
	"github.com/open-policy-agent/opa/ast"
)

// List of builtins supported by OPA Wasm compiler - what is not on this list has
// to be provided by the host
//
// This map is a copy and paste from https://github.com/open-policy-agent/opa/blob/251d4adce279a9f0613c5fd08731827ec72c8937/internal/compiler/wasm/wasm.go#L81
// Unfortunately this isn't an exported object.
// We only care about the keys of this map, we keep the map around
// to simplify the process of keeping this data up-to-date
var builtinsFunctions = map[string]string{
	ast.Plus.Name:                       "opa_arith_plus",
	ast.Minus.Name:                      "opa_arith_minus",
	ast.Multiply.Name:                   "opa_arith_multiply",
	ast.Divide.Name:                     "opa_arith_divide",
	ast.Abs.Name:                        "opa_arith_abs",
	ast.Round.Name:                      "opa_arith_round",
	ast.Ceil.Name:                       "opa_arith_ceil",
	ast.Floor.Name:                      "opa_arith_floor",
	ast.Rem.Name:                        "opa_arith_rem",
	ast.ArrayConcat.Name:                "opa_array_concat",
	ast.ArraySlice.Name:                 "opa_array_slice",
	ast.SetDiff.Name:                    "opa_set_diff",
	ast.And.Name:                        "opa_set_intersection",
	ast.Or.Name:                         "opa_set_union",
	ast.Intersection.Name:               "opa_sets_intersection",
	ast.Union.Name:                      "opa_sets_union",
	ast.IsNumber.Name:                   "opa_types_is_number",
	ast.IsString.Name:                   "opa_types_is_string",
	ast.IsBoolean.Name:                  "opa_types_is_boolean",
	ast.IsArray.Name:                    "opa_types_is_array",
	ast.IsSet.Name:                      "opa_types_is_set",
	ast.IsObject.Name:                   "opa_types_is_object",
	ast.IsNull.Name:                     "opa_types_is_null",
	ast.TypeNameBuiltin.Name:            "opa_types_name",
	ast.BitsOr.Name:                     "opa_bits_or",
	ast.BitsAnd.Name:                    "opa_bits_and",
	ast.BitsNegate.Name:                 "opa_bits_negate",
	ast.BitsXOr.Name:                    "opa_bits_xor",
	ast.BitsShiftLeft.Name:              "opa_bits_shiftleft",
	ast.BitsShiftRight.Name:             "opa_bits_shiftright",
	ast.Count.Name:                      "opa_agg_count",
	ast.Sum.Name:                        "opa_agg_sum",
	ast.Product.Name:                    "opa_agg_product",
	ast.Max.Name:                        "opa_agg_max",
	ast.Min.Name:                        "opa_agg_min",
	ast.Sort.Name:                       "opa_agg_sort",
	ast.All.Name:                        "opa_agg_all",
	ast.Any.Name:                        "opa_agg_any",
	ast.Base64IsValid.Name:              "opa_base64_is_valid",
	ast.Base64Decode.Name:               "opa_base64_decode",
	ast.Base64Encode.Name:               "opa_base64_encode",
	ast.Base64UrlEncode.Name:            "opa_base64_url_encode",
	ast.Base64UrlDecode.Name:            "opa_base64_url_decode",
	ast.NetCIDRContains.Name:            "opa_cidr_contains",
	ast.NetCIDROverlap.Name:             "opa_cidr_contains",
	ast.NetCIDRIntersects.Name:          "opa_cidr_intersects",
	ast.Equal.Name:                      "opa_cmp_eq",
	ast.GreaterThan.Name:                "opa_cmp_gt",
	ast.GreaterThanEq.Name:              "opa_cmp_gte",
	ast.LessThan.Name:                   "opa_cmp_lt",
	ast.LessThanEq.Name:                 "opa_cmp_lte",
	ast.NotEqual.Name:                   "opa_cmp_neq",
	ast.GlobMatch.Name:                  "opa_glob_match",
	ast.JSONMarshal.Name:                "opa_json_marshal",
	ast.JSONUnmarshal.Name:              "opa_json_unmarshal",
	ast.ObjectFilter.Name:               "builtin_object_filter",
	ast.ObjectGet.Name:                  "builtin_object_get",
	ast.ObjectRemove.Name:               "builtin_object_remove",
	ast.ObjectUnion.Name:                "builtin_object_union",
	ast.Concat.Name:                     "opa_strings_concat",
	ast.FormatInt.Name:                  "opa_strings_format_int",
	ast.IndexOf.Name:                    "opa_strings_indexof",
	ast.Substring.Name:                  "opa_strings_substring",
	ast.Lower.Name:                      "opa_strings_lower",
	ast.Upper.Name:                      "opa_strings_upper",
	ast.Contains.Name:                   "opa_strings_contains",
	ast.StartsWith.Name:                 "opa_strings_startswith",
	ast.EndsWith.Name:                   "opa_strings_endswith",
	ast.Split.Name:                      "opa_strings_split",
	ast.Replace.Name:                    "opa_strings_replace",
	ast.ReplaceN.Name:                   "opa_strings_replace_n",
	ast.Trim.Name:                       "opa_strings_trim",
	ast.TrimLeft.Name:                   "opa_strings_trim_left",
	ast.TrimPrefix.Name:                 "opa_strings_trim_prefix",
	ast.TrimRight.Name:                  "opa_strings_trim_right",
	ast.TrimSuffix.Name:                 "opa_strings_trim_suffix",
	ast.TrimSpace.Name:                  "opa_strings_trim_space",
	ast.NumbersRange.Name:               "opa_numbers_range",
	ast.ToNumber.Name:                   "opa_to_number",
	ast.WalkBuiltin.Name:                "opa_value_transitive_closure",
	ast.ReachableBuiltin.Name:           "builtin_graph_reachable",
	ast.RegexIsValid.Name:               "opa_regex_is_valid",
	ast.RegexMatch.Name:                 "opa_regex_match",
	ast.RegexMatchDeprecated.Name:       "opa_regex_match",
	ast.RegexFindAllStringSubmatch.Name: "opa_regex_find_all_string_submatch",
	ast.JSONRemove.Name:                 "builtin_json_remove",
	ast.JSONFilter.Name:                 "builtin_json_filter",
}

// Builtins returns a Set with the names of all the builtins that are
// provided by Wasm policies
func Builtins() mapset.Set {
	wasmBuiltins := mapset.NewSet()
	for k := range builtinsFunctions {
		wasmBuiltins.Add(string(k))
	}

	return wasmBuiltins
}
