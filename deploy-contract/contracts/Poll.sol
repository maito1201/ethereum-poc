//SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.0;

contract Poll {
    event Voted(address _voter, uint _value);
    mapping(address => uint) public votes;
    string pollSubject = "this is test poll";

    constructor() {
        getPoll();
    }

    function getPoll() public view returns (string memory) {
        return pollSubject;
    }

    function vote(uint selection) public {
        emit Voted(msg.sender, selection);
        require (votes[msg.sender] == 0);
        require (selection > 0 && selection < 3);
        votes[msg.sender] = selection;
    }

    function checkPoll() public view returns (uint) {
        return votes[msg.sender];
    }
}
