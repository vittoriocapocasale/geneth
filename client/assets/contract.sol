// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Counter {

  uint64 count;

  
  function increment () public {
    count=count+1;
  }   
  
  function get() public view returns (uint64){
    return (count);
  }
}
