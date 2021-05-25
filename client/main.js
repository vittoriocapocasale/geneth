'use strict';

const utils = require('./utils')

main();


async function main(){

    var options = process.argv.slice(2);
    switch (options[0]) {
      case 'read':
        await utils.getStoredObject(utils.loadSettings('./assets/environmentNode1.json'));
        break;
      case 'deploy':
        let env = utils.loadSettings('./assets/environmentNode0.json');
        let web3= utils.getWebHandle(env.url);
        await utils.deployContract(web3, env);
        break;
      case 'increment':
        await utils.increment(utils.loadSettings('./assets/environmentNode0.json'));
        break;
      default:
        console.log('Unknown...quitting');
    }

}

