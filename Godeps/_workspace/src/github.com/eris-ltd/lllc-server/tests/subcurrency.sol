
contract Subcurrency {
	address creator;
	mapping (address => uint) balances;
	function Subcurrency() {
		creator = msg.sender;
	}

	function mint(address owner, uint amt){
		if (tx.origin != creator) return;
		balances[owner] += amt;
	}

	// sends from the origin to the caller
	function transfer(uint amt){
		if (balances[tx.origin] < amt) return;
		balances[tx.origin] -= amt;
		balances[msg.sender] += amt;
	}

	function send(address rec, uint amt){
		if (balances[msg.sender] < amt) return;
		balances[msg.sender] -= amt;
		balances[rec] += amt;
	}

	function query(address addr) constant returns (uint balance) {
		return balances[addr];
	}
}

