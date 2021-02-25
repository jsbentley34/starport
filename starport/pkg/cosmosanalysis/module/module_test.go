package module

import (
	"github.com/kr/pretty"
)

func ExampleDiscover() {
	pretty.Println(Discover("/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon"))
	// outputs:
	// []module.Module{
	//   {
	//     Name:            "moon",
	//     TypesImportPath: "github.com/test/moon/x/moon/types",
	//     Msgs: {
	//       { Name:"MsgCreateUser", URI:"test.moon.moon.MsgCreateUser", FilePath:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/tx.proto" },
	//       { Name:"MsgDeleteUser", URI:"test.moon.moon.MsgDeleteUser", FilePath:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/tx.proto" },
	//       { Name:"MsgUpdateUser", URI:"test.moon.moon.MsgUpdateUser", FilePath:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/tx.proto" },
	//     },
	//   },
	// } nil

}
