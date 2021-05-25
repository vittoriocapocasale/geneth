'use strict';

const fs = require('fs');
const Web3 = require('web3');

module.exports = {
  loadSettings:function (filepath){
    let rawdata = fs.readFileSync(filepath);
    let settings = JSON.parse(rawdata);
    return settings
  },
  getWebHandle: function(url){
    console.log(url)
    const web3 = new Web3(new Web3.providers.WebsocketProvider(url));
    return web3
  },
  getContract: function(web3){
    let abi =JSON.parse(fs.readFileSync('./assets/build/contract_sol_Counter.abi'));
    let address = fs.readFileSync('./assets/address.txt').toString();
    let contract = new web3.eth.Contract(abi, address);
    return contract
  },
  getStoredObject: async function(env) {
        let web3 = this.getWebHandle(env.url)
        let contract = this.getContract(web3)
        let result = contract.methods.get()
        .call({
           from: env.accAddr[0],//[this.submitted%this.env.accAddr.length],
           gasPrice: await web3.eth.getGasPrice(), 
           gas: 800000,
           })
        let r = await result;
        console.log("Current count:",r);
        web3.currentProvider.disconnect();
  },
  increment: async function(env) {
        let web3 = this.getWebHandle(env.url)
        let contract = this.getContract(web3)
        let result = contract.methods.increment()
        .send({
           from: env.accAddr[0],//[this.submitted%this.env.accAddr.length],
           gasPrice: await web3.eth.getGasPrice(), 
           gas: 800000,
           })
        let r = await result;
        console.log(r);
        web3.currentProvider.disconnect();
  },
  fileAppend: async function(str){
    var stream = fs.createWriteStream("./assets/results.txt", {flags:'a'});
    stream.write(str+',\n')
  },
  deployContract: async function (web3, env) {
    let bytecode = fs.readFileSync('./assets/build/contract_sol_Counter.bin');
    let abi =JSON.parse(fs.readFileSync('./assets/build/contract_sol_Counter.abi'));
    let contract = new web3.eth.Contract(abi, null, {data: '0x'+bytecode, from: env.accAddr[0]});
    let result;
    try{
        result = await contract.deploy().send()
    }
    catch(e){console.log(e); result=null}
    fs.writeFile('./assets/address.txt', result.options.address , function (err) {
      if (err) console.log(err);
    });
    console.log(result.options.address);
    web3.currentProvider.disconnect();
  }
};




  
