package test

type ByteAlias byte

type BoolAlias bool

type Int32Alias int32

type Float32Alias float32

type StringAlias string

type SliceStringAlias []string

type SlicePtrStringAlias []*string

type MapStringStringAlias map[string]string

type Inner struct {
	Byte            byte
	BytePtr         *byte
	ByteAlias       ByteAlias
	ByteAliasPtr    *ByteAlias
	Bool            bool
	BoolPtr         *bool
	BoolAlias       BoolAlias
	BoolAliasPtr    *BoolAlias
	Int8            int8
	Int8Ptr         *int8
	Int16           int16
	Int16Ptr        *int16
	Int32           int32
	Int32Ptr        *int32
	Int32Alias      Int32Alias
	Int32AliasPtr   *Int32Alias
	Uint8           uint8
	Uint8Ptr        *uint8
	Uint16          uint16
	Uint16Ptr       *uint16
	Uint32          uint32
	Uint32Ptr       *uint32
	Float32         float32
	Float32Ptr      *float32
	Float32Alias    Float32Alias
	Float32AliasPtr *Float32Alias
	Float64         float64
	Float64Ptr      *float64
	String          string
	StringPtr       *string
	StringAlias     StringAlias
	StringAliasPtr  *StringAlias
	Struct          struct {
		Byte                byte
		BytePtr             *byte
		ByteAlias           ByteAlias
		ByteAliasPtr        *ByteAlias
		Bool                bool
		BoolPtr             *bool
		BoolAlias           BoolAlias
		BoolAliasPtr        *BoolAlias
		Int8                int8
		Int8Ptr             *int8
		Int16               int16
		Int16Ptr            *int16
		Int32               int32
		Int32Ptr            *int32
		Int32Alias          Int32Alias
		Int32AliasPtr       *Int32Alias
		Uint8               uint8
		Uint8Ptr            *uint8
		Uint16              uint16
		Uint16Ptr           *uint16
		Uint32              uint32
		Uint32Ptr           *uint32
		Float32             float32
		Float32Ptr          *float32
		Float32Alias        Float32Alias
		Float32AliasPtr     *Float32Alias
		Float64             float64
		Float64Ptr          *float64
		String              string
		StringPtr           *string
		StringAlias         StringAlias
		StringAliasPtr      *StringAlias
		Struct              struct{}
		StructPtr           *Inner
		SliceString         []string
		SliceStringAlias    SliceStringAlias
		SlicePtrString      []*string
		SliceStringPtrAlias SlicePtrStringAlias
		SliceStringPtr      *[]string
		SliceByte           []byte
	}
	StructPtr *struct {
		Byte                    byte
		BytePtr                 *byte
		ByteAlias               ByteAlias
		ByteAliasPtr            *ByteAlias
		Bool                    bool
		BoolPtr                 *bool
		BoolAlias               BoolAlias
		BoolAliasPtr            *BoolAlias
		Int8                    int8
		Int8Ptr                 *int8
		Int16                   int16
		Int16Ptr                *int16
		Int32                   int32
		Int32Ptr                *int32
		Int32Alias              Int32Alias
		Int32AliasPtr           *Int32Alias
		Uint8                   uint8
		Uint8Ptr                *uint8
		Uint16                  uint16
		Uint16Ptr               *uint16
		Uint32                  uint32
		Uint32Ptr               *uint32
		Float32                 float32
		Float32Ptr              *float32
		Float32Alias            Float32Alias
		Float32AliasPtr         *Float32Alias
		Float64                 float64
		Float64Ptr              *float64
		String                  string
		StringPtr               *string
		StringAlias             StringAlias
		StringAliasPtr          *StringAlias
		Struct                  struct{}
		StructPtr               *Inner
		SliceString             []string
		SliceStringAlias        SliceStringAlias
		SlicePtrString          []*string
		SliceStringPtrAlias     SlicePtrStringAlias
		SliceStringPtr          *[]string
		SliceByte               []byte
		MapStringString         map[string]string
		MapStringStringPtr      *map[string]string
		MapStringStringAlias    MapStringStringAlias
		MapStringStringAliasPtr *MapStringStringAlias
	}
	SliceString             []string
	SliceStringAlias        SliceStringAlias
	SlicePtrString          []*string
	SliceStringPtrAlias     SlicePtrStringAlias
	SliceStringPtr          *[]string
	SliceByte               []byte
	MapStringString         map[string]string
	MapStringStringPtr      *map[string]string
	MapStringStringAlias    MapStringStringAlias
	MapStringStringAliasPtr *MapStringStringAlias
}

type T struct {
	Byte            byte
	BytePtr         *byte
	ByteAlias       ByteAlias
	ByteAliasPtr    *ByteAlias
	Bool            bool
	BoolPtr         *bool
	BoolAlias       BoolAlias
	BoolAliasPtr    *BoolAlias
	Int8            int8
	Int8Ptr         *int8
	Int16           int16
	Int16Ptr        *int16
	Int32           int32
	Int32Ptr        *int32
	Int32Alias      Int32Alias
	Int32AliasPtr   *Int32Alias
	Uint8           uint8
	Uint8Ptr        *uint8
	Uint16          uint16
	Uint16Ptr       *uint16
	Uint32          uint32
	Uint32Ptr       *uint32
	Float32         float32
	Float32Ptr      *float32
	Float32Alias    Float32Alias
	Float32AliasPtr *Float32Alias
	Float64         float64
	Float64Ptr      *float64
	String          string
	StringPtr       *string
	StringAlias     StringAlias
	StringAliasPtr  *StringAlias
	StructPtr       *Inner
	Struct          struct {
		Struct struct {
			Struct struct {
				Struct struct {
					String string
				}
			}
		}
	}
	SliceString             []string
	SliceStringAlias        SliceStringAlias
	SlicePtrString          []*string
	SliceStringPtrAlias     SlicePtrStringAlias
	SliceStringPtr          *[]string
	MapStringString         map[string]string
	MapStringStringPtr      *map[string]string
	MapStringStringAlias    MapStringStringAlias
	MapStringStringAliasPtr *MapStringStringAlias
}
