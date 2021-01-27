package goja

import (
	"testing"

	gocmp "github.com/google/go-cmp/cmp"
)

func TestObject_(t *testing.T) {
	vm := New()

	for _, tt := range []struct {
		js       string
		expected interface{}
	}{
		{
			`
            var abc = Object.getOwnPropertyDescriptor(Object, "prototype");
            [ [ typeof Object.prototype, abc.writable, abc.enumerable, abc.configurable ],
            ];
		`, "object,false,false,false",
		},
	} {
		actual, err := vm.RunString(tt.js)
		if err != nil {
			t.Fatal(err)
		}
		if failed := gocmp.Diff(actual.String(), tt.expected); failed != "" {
			t.Fatal(failed)
		}
	}
}

func TestObject_new(t *testing.T) {
	vm := New()

	for _, tt := range []struct {
		js       string
		expected interface{}
	}{
		{
			`
            [ new Object("abc"), new Object(2+2) ];
        `, "abc,4",
		},
	} {
		actual, err := vm.RunString(tt.js)
		if err != nil {
			t.Fatal(err)
		}
		if failed := gocmp.Diff(actual.String(), tt.expected); failed != "" {
			t.Fatal(failed)
		}
	}
}

func TestObject_keys(t *testing.T) {
	vm := New()

	for _, tt := range []struct {
		js       string
		expected interface{}
	}{
		{
			`Object.keys({ abc:undefined, def:undefined })`,
			"abc,def",
		},
		{
			`
            function abc() {
                this.abc = undefined;
                this.def = undefined;
            }
            Object.keys(new abc())
		`,
			"abc,def",
		},
		{
			`
	        function def() {
	            this.ghi = undefined;
	        }
	        def.prototype = new abc();
	        Object.keys(new def());
		`,
			"ghi",
		},
		{
			` (function(abc, def, ghi){
	return Object.keys(arguments)
})(undefined, undefined);
`,
			"0,1",
		},
		{
			`
	        (function(abc, def, ghi){
	            return Object.keys(arguments)
	        })(undefined, undefined, undefined, undefined);
		`,
			"0,1,2,3",
		},
	} {
		actual, err := vm.RunString(tt.js)
		if err != nil {
			t.Fatal(err)
		}
		if failed := gocmp.Diff(actual.String(), tt.expected); failed != "" {
			t.Fatal(failed)
		}
	}
}

func TestObject_values(t *testing.T) {
	vm := New()

	for _, tt := range []struct {
		js       string
		expected interface{}
	}{
		{
			`Object.values({ k1: 'abc', k2 :'def' })`, "abc,def",
		},

		{
			`
						function abc() {
							this.k1 = "abc";
							this.k2 = "def";
						}
						Object.values(new abc());
				`, "abc,def",
		},

		{
			`
						function def() {
							this.k3 = "ghi";
						}
						def.prototype = new abc();
						Object.values(new def());
				`, "ghi",
		},

		{
			`
						var ghi = Object.create(
                {
                    k1: "abc",
                    k2: "def"
                },
                {
                    k3: { value: "ghi", enumerable: true },
                    k4: { value: "jkl", enumerable: false }
                }
            );
            Object.values(ghi);
				`, "ghi",
		},

		{
			`
            (function(abc, def, ghi){
                return Object.values(arguments)
            })(0, 1);
        `, "0,1",
		},

		{
			`
            (function(abc, def, ghi){
                return Object.values(arguments)
            })(0, 1, 2, 3);
        `, "0,1,2,3",
		},
	} {
		actual, err := vm.RunString(tt.js)
		if err != nil {
			t.Fatal(err)
		}
		if failed := gocmp.Diff(actual.String(), tt.expected); failed != "" {
			t.Fatal(failed)
		}
	}
}

func TestObject_entries(t *testing.T) {
	vm := New()

	for _, tt := range []struct {
		js       string
		expected interface{}
	}{
		{
			`Object.entries({ k1: 'abc', k2 :'def' })`, "k1,abc,k2,def",
		},

		{
			`
 			      var e = Object.entries({ k1: 'abc', k2 :'def' });
 						[ e[0][0], e[0][1], e[1][0], e[1][1], ];
 				 `, "k1,abc,k2,def",
		},
		{
			`
 						function abc() {
 							this.k1 = "abc";
 							this.k2 = "def";
 						}
 						Object.entries(new abc());
 				`, "k1,abc,k2,def",
		},
		{
			`
 						function def() {
 							this.k3 = "ghi";
 						}
 						def.prototype = new abc();
 						Object.entries(new def());
				 `,
			"k3,ghi",
		},
		{
			`
 						var ghi = Object.create(
 		            {
 		                k1: "abc",
 		                k2: "def"
 		            },
 		            {
 		                k3: { value: "ghi", enumerable: true },
 		                k4: { value: "jkl", enumerable: false }
 		            }
 		        );
 		        Object.entries(ghi);
				 `,
			"k3,ghi",
		},
		{
			`
 		        (function(abc, def, ghi){
 		            return Object.entries(arguments)
 		        })(0, 1);
 		    `, "0,0,1,1",
		},
		{
			`
 		        (function(abc, def, ghi){
 		            return Object.entries(arguments)
 		        })(0, 1, 2, 3);
 		    `, "0,0,1,1,2,2,3,3",
		},
	} {
		actual, err := vm.RunString(tt.js)
		if err != nil {
			t.Fatal(err)
		}
		if failed := gocmp.Diff(actual.String(), tt.expected); failed != "" {
			t.Fatal(failed)
		}
	}
}

func TestObject_fromEntries(t *testing.T) {
	vm := New()

	for _, tt := range []struct {
		js       string
		expected interface{}
	}{
		{
			`
 					 var o = Object.fromEntries([['a', 1], ['b', true], ['c', 'sea']]);
 					 [ o.a, o.b, o.c ]
 				 `, "1,true,sea",
		},
	} {
		actual, err := vm.RunString(tt.js)
		if err != nil {
			t.Fatal(err)
		}
		if failed := gocmp.Diff(actual.String(), tt.expected); failed != "" {
			t.Fatal(failed)
		}
	}
}

func TestObject_defineSetter(t *testing.T) {
	vm := New()

	src := `
		var o = {};
		o.__defineSetter__('value', function(val) { this.anotherValue = val; });
		o.value = 5;
		[o.value + "," + o.anotherValue]
	`

	actual, err := vm.RunString(src)
	if err != nil {
		t.Fatal(err)
	}
	if failed := gocmp.Diff(actual.String(), "undefined,5"); failed != "" {
		t.Fatal(failed)
	}
}

func TestObject_defineSetterError(t *testing.T) {
	vm := New()

	for _, tc := range []struct {
		desc string
		src  string
	}{
		{
			"with an undefined function",
			`
			var o = {};
			Object.prototype.__defineSetter__('value');
		`,
		},
		{
			"with an empty object",
			`
				var o = {};
				Object.prototype.__defineSetter__('value', {});
			`,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := vm.RunString(tc.src)
			if err == nil {
				t.Fatal("Expected an error")
			}

			if err.Error() != "TypeError: Object.prototype.__defineSetter__: Expecting function" {
				t.Fatal("Unexpected error: ", err)
			}
		})
	}
}

func TestObject_toString(t *testing.T) {
	vm := New()

	for _, tc := range []struct {
		desc     string
		src      string
		expected string
	}{
		{
			"with a Date",
			"Object.prototype.toString.call(new Date);",
			"[object Date]",
		},
		{
			"with a String",
			"Object.prototype.toString.call(new String);",
			"[object String]",
		},
		{
			"with Math",
			"Object.prototype.toString.call(Math);",
			"[object Math]",
		},
		{
			"with undefined",
			"Object.prototype.toString.call(undefined);",
			"[object Undefined]",
		},
		{
			"with null",
			"Object.prototype.toString.call(null);",
			"[object Null]",
		},
		{
			"with an object",
			"Object.prototype.toString.call({});",
			"[object Object]",
		},
		{
			"with a number",
			"Object.prototype.toString.call(3);",
			"[object Number]",
		},
		{
			"with an array",
			"Object.prototype.toString.call([1]);",
			"[object Array]",
		},
		{
			"with a toStringTag property",
			`
				var myDate = new Date();
				myDate[Symbol.toStringTag] = 'myDate';
				Object.prototype.toString.call(myDate);
			`,
			"[object myDate]",
		},
		{
			"with a named function",
			`
				var myFunc = function hello(){};
				Object.prototype.toString.call(myFunc);
			`,
			"[object Function]",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			actual, err := vm.RunString(tc.src)
			if err != nil {
				t.Fatal("Unexpected error: ", err)
			}
			if failed := gocmp.Diff(actual.String(), tc.expected); failed != "" {
				t.Fatal(failed)
			}
		})
	}
}