package encoder

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alphabill-org/alphabill-go-base/types"
)

func Test_TXSystemEncoder_trigger(t *testing.T) {
	/*
		If test here fails it's probably because some data structure (or rather how it's
		serialized for Rust SDK) has been changed without versioning?
		Also, the Rust predicates SDK likely needs to be updated! See the
		other tests here to generate tests for Rust SDK.
	*/

	// encoder is stateless, can be shared between tests
	enc, err := New()
	require.NoError(t, err)

	t.Run("txRecord", func(t *testing.T) {
		// ver 1 of the txRec just contains txo handle so no need to fill out all the fields...
		// getHandle is called once, always return 1 as the handle
		getHandle := func(obj any) uint64 { return 1 }
		tx, err := (&types.TransactionOrder{Version: 1}).MarshalCBOR()
		require.NoError(t, err)
		buf, err := enc.Encode(&types.TransactionRecord{Version: 1, TransactionOrder: tx}, 1, getHandle)
		require.NoError(t, err)
		require.Equal(t, []byte{0x1, 0x2, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, buf)
	})

	t.Run("txOrder", func(t *testing.T) {
		getHandle := func(obj any) uint64 { t.Errorf("unexpected call of getHandle(%T)", obj); return 0 }
		// ver 1 of the txOrder
		txo := &types.TransactionOrder{
			Version: 1,
			Payload: types.Payload{
				PartitionID: 7,
				Type:        22,
				UnitID:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				ClientMetadata: &types.ClientMetadata{
					ReferenceNumber: []byte("ref-no"),
				},
			},
		}
		buf, err := enc.Encode(txo, 1, getHandle)
		require.NoError(t, err)
		require.Equal(t, []byte{0x1, 0x4, 0x0, 0x0, 0x2, 0x3, 0x7, 0x0, 0x0, 0x0, 0x3, 0x1, 0xa, 0x0, 0x0, 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0x4, 0x4, 0x16, 0x0, 0x5, 0x1, 0x6, 0x0, 0x0, 0x0, 0x72, 0x65, 0x66, 0x2d, 0x6e, 0x6f}, buf)
	})

	t.Run("byte slice", func(t *testing.T) {
		// byte slice is returned exactly as-is
		buf, err := enc.Encode([]byte{0, 1, 127, 128, 255}, 1, nil)
		require.NoError(t, err)
		require.Equal(t, []byte{0, 1, 127, 128, 255}, buf)
	})

	t.Run("types.RawCBOR", func(t *testing.T) {
		// types.RawCBOR is returned exactly as-is but as byte slice (ie type cast)
		buf, err := enc.Encode(types.RawCBOR{0, 1, 127, 128, 255}, 1, nil)
		require.NoError(t, err)
		require.Equal(t, []byte{0, 1, 127, 128, 255}, buf)
	})
}

func Test_generate_TXSTestsData(t *testing.T) {
	t.Skip("generate test data for Rust predicate SDK")
	/*
		Generate inputs for tests loading generic tx system objects.
		In the Rust SDK tests are in the file "src/txsystem.rs"
		We only generate input for current struct version here, maintain
		tests for relevant versions in the Rust code!

		To use: comment out the Skip statement in the beginning, run the
		test(s) to generate data, copy it to Rust project and uncomment
		the Skip here.
	*/

	var hid atomic.Uint64
	getHandle := func(obj any) uint64 { return hid.Add(1) }
	// encoder is stateless, can be shared between tests
	enc, err := New()
	require.NoError(t, err)

	t.Run("txRecord", func(t *testing.T) {
		// ver 1 of the txRec just contains txo handle so no need to fill out all the fields...
		tx, err := (&types.TransactionOrder{Version: 1}).MarshalCBOR()
		require.NoError(t, err)
		buf, err := enc.Encode(&types.TransactionRecord{Version: 1, TransactionOrder: tx}, 1, getHandle)
		require.NoError(t, err)
		t.Errorf("\nlet data: &mut [u8] = &mut [%s];", bytesAsHex(t, buf))
	})

	t.Run("txOrder", func(t *testing.T) {
		// ver 1 of the txOrder
		txo := &types.TransactionOrder{
			Version: 1,
			Payload: types.Payload{
				PartitionID: 7,
				Type:        22,
				UnitID:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				ClientMetadata: &types.ClientMetadata{
					ReferenceNumber: []byte("ref-no"),
				},
			},
		}
		buf, err := enc.Encode(txo, 1, getHandle)
		require.NoError(t, err)
		t.Errorf("\n// %#v\n// %#v\nlet data: &mut [u8] = &mut [%s];", txo.Payload, txo.ClientMetadata, bytesAsHex(t, buf))
	})
}

func Test_generateDecoderTests(t *testing.T) {
	t.Skip("generate test data for Rust predicate SDK Decoder")
	/*
		To test that the Decoder in Rust SDK is able to decode data generated by the host
		we generate these Rust tests:
		- set the (absolute) path in the os.Create statement below to where to save the
		  generated code;
		- comment out the Skip statement in the beginning;
		- run the test(s);
		- copy the generated code into src/decoder.rs in Rust SDK;
		- NB! the test fails on purpose to remind you to enable the Skip statement again!

		At this stage doing all this manual work seems to be acceptable (rather than building
		some test harness to do it)...
	*/

	fOut, err := os.Create("decoder_test.rs")
	require.NoError(t, err)
	defer fOut.Close()

	type encValue struct {
		tag   uint8
		value any
	}

	t.Run("Decode Value", func(t *testing.T) {
		// test decoding different types as Value enum defined in the Rust SDK
		values := []encValue{
			{value: uint32(101)},
			{value: uint64(64)},
			{value: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{value: "this is string"},
			{value: []any{uint32(32), uint64(64), "AB"}},
			// zero values
			{value: ""},
			{value: []byte(nil)},
			{value: []any{}},
		}
		out := bytes.NewBufferString("\n#[test]\nfn decode_value() {")
		out.WriteString("\n// test generated by Go backend!")
		for _, v := range values {
			enc := TVEnc{}
			enc.Encode(v.value)
			buf, err := enc.Bytes()
			require.NoError(t, err)
			out.WriteString(fmt.Sprintf("\n// Go value %T = %#[1]v\n", v.value))
			out.WriteString("let data: &mut [u8] = &mut [")
			out.WriteString(bytesAsHex(t, buf))
			out.WriteString("];\n")
			out.WriteString("let mut dec = Decoder::new(data);\n")
			out.WriteString(fmt.Sprintf("assert_eq!(dec.value(), %s);\n", rustValue(t, v.value)))
		}
		out.WriteString("}\n")
		_, err = out.WriteTo(fOut)
		require.NoError(t, err)
	})

	t.Run("TagValueIter", func(t *testing.T) {
		// test for parsing data sent from host using TagValueIter
		values := []encValue{
			{tag: 1, value: uint32(0xff00ff00)},
			{tag: 4, value: uint64(0xff00ff00)},
			{tag: 2, value: "token"},
		}
		out := bytes.NewBufferString("\n#[test]\nfn iterator() {")
		out.WriteString("\n// test generated by Go backend!\n")
		enc := TVEnc{}
		rv := []string{}
		for _, v := range values {
			enc.EncodeTagged(v.tag, v.value)
			rv = append(rv, fmt.Sprintf("(%d, %s)", v.tag, rustValue(t, v.value)))
		}
		buf, err := enc.Bytes()
		require.NoError(t, err)
		out.WriteString("let data: &mut [u8] = &mut [")
		out.WriteString(bytesAsHex(t, buf))
		out.WriteString("];\n")
		out.WriteString("let dec = TagValueIter::new(&data);\nlet items: Vec<(u8, Value)> = dec.collect();\n")
		out.WriteString(fmt.Sprintf("assert_eq!(items, vec![%s]);\n", strings.Join(rv, ", ")))
		out.WriteString("}\n")
		_, err = out.WriteTo(fOut)
		require.NoError(t, err)
	})

	t.Error("Do NOT remove this error - instead add/enable the t.Skip statement in the beginning of the test! " +
		"This error is here so that the test is enabled only for re-generating the test code for Rust SDK and then disabled again.")
}

// v as Rust SDK Value enum
func rustValue(t *testing.T, v any) string {
	t.Helper()
	switch tv := v.(type) {
	case uint32:
		return fmt.Sprintf("Value::U32(%d)", tv)
	case uint64:
		return fmt.Sprintf("Value::U64(%d)", tv)
	case []byte:
		return fmt.Sprintf("Value::Bytes(vec![%v])", bytesAsHex(t, tv))
	case string:
		return fmt.Sprintf("Value::String(%q.to_string())", tv)
	case []any:
		out := "Value::Array(vec!["
		for x, v := range tv {
			if x > 0 {
				out += ", "
			}
			out += rustValue(t, v)
		}
		return out + "])"
	default:
		t.Errorf("unsupported type %T", v)
		return ""
	}
}

// byte slice in Rust syntax (without enclosing [])
func bytesAsHex(t *testing.T, b []byte) string {
	t.Helper()
	out := bytes.Buffer{}
	for _, v := range b {
		out.WriteString(fmt.Sprintf("0x%x, ", v))
	}
	return strings.TrimSuffix(out.String(), ", ")
}
