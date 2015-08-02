
contract AtomicSwap {
	bool isOpen; // if the open function ran successfuly and hasnt been closed
	address currency; // currency contract for the coin we're swapping
	uint256 hashLock; // hash of the secret
	uint swapAmt; // amount to send
	address sender; // she who opens the channel
	address receiver; // she who receives on the channel
	uint startHeight; // block height at which channel opened
	uint lockTime; // how many blocks the funds are locked for before refund

	// open the channel by specifying all the params and the current height
	function open(address currencyContract, uint amt, address rec, uint256 hash, uint lockT){
		if (isOpen) return; // return if a swap is already in progress
		currency = currencyContract;
		sender = msg.sender;

		// check the senders balance in subcurrency contract
		uint bal = uint(currency.call(bytes4(sha3("query(address)")), sender));
		if (bal < amt) return; // not enough funds
		receiver = rec;

		// transfer from sender to this contract
		currency.call(bytes4(sha3("transfer(uint)")), amt);

		// set final params
		swapAmt = amt;
		hashLock = hash;
		startHeight = block.number;
		lockTime = lockT;

		// swap is now active!
		isOpen = true;
	}

	function refund(){
		if (!isOpen) return;
		if (block.number > startHeight + lockTime){
			currency.call(bytes4(sha3("send(address,uint)")), sender,swapAmt);
			isOpen = false;
		}
	}

	function redeem(uint256 secret){
		if (!isOpen) return;
		if (sha3(secret) == hashLock){
			currency.call(bytes4(sha3("send(address,uint)")), receiver,swapAmt);
		}
	}
}

