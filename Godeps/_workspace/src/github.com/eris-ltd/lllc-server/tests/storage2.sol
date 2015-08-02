contract SimpleStorage {
	struct permdata {
		uint permnum;
		uint globalval;
	}
	struct userdata {
		mapping (uint => uint) permvals;
		mapping (uint => bool) exclude;
	}

	mapping (address => userdata) users;
	mapping (bytes32 => permdata) perms;
	uint pcount;
	mapping (bytes32 => address) names;

	function SimpleStorage() {
		pcount =  1;
		permdata p = perms["DOUG"];
		userdata u = users[tx.origin];
		p.permnum = pcount;
		p.globalval = 0;

		//Give creator DOUG permissions
		u.permvals[pcount] = 1;
	}

	//NameReg Functionality
	function checkName(bytes32 name) returns (address ret){
		return names[name];
	}

	function register(bytes32 name, address addr) returns (uint ret) {
		//Check for doug perms for ORIGIN
		if (getPerm(tx.origin, "DOUG")!=1) return 0;

		names[name] = addr;
		return 1;
	}

	//Permission Functionality
	function getPerm(address addr, bytes32 permName) returns (uint ret) {
		permdata p = perms[permName];
		userdata u = users[addr];

		//Check if permName even exists
		if (p.permnum == 0) return 0;

		//Process perm
		if (u.permvals[p.permnum] == 0){
			if (u.exclude[p.permnum]) return u.permvals[p.permnum];
			else return p.globalval;
		}
		else { 
			return u.permvals[p.permnum]; 
		}
	}


    uint storedData;
    uint otherThing;
    function set(uint x, uint y) {
        storedData = x;
        otherThing = y;
    }
    function getX() constant returns (uint retVal) {
        return storedData;
    }

    function getY() constant returns (uint retVal) {
        return storedData;
    }

    function doSomething() {
	uint a = 3 + 4;
	uint b = 5 + 6;
	uint c = a - b;
    }

    function set2(uint x, uint y) {
        storedData = x;
        otherThing = y;
    }

    function set3(uint x, uint y) {
        storedData = x;
        otherThing = y;
    }

    function set4(uint x, uint y) {
        storedData = x;
        otherThing = y;
    }

    function set5(uint x, uint y) {
        storedData = x;
        otherThing = y;
    }
}

