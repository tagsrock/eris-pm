 contract multiReturn {

	struct Strings {
		string filename;
		string username;
	}

	function getStatus() returns (uint, uint) {
  		return (1, 2);
	}
	function getStrings() returns (string filename, string username) {
	    
		return (filename, username);
	}
}

