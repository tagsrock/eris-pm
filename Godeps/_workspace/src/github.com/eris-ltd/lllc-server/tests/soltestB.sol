contract B{
	function B(){
	}

	function send(address anA){
		anA.call(bytes4(sha3("hi()")))	;
	}
}
