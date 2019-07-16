const Fabric_Client = require("fabric-client");
const path = require("path");
const util = require("util");
const os = require("os");

const fabric_client = new Fabric_Client();

// setup the fabric network
const channel = fabric_client.newChannel("jewelchannel");

const peer = fabric_client.newPeer("grpc://0.0.0.0:7051");
channel.addPeer(peer);

const order = fabric_client.newOrderer("grpc://0.0.0.0:7050");
channel.addOrderer(order);

let member_user = null;
const store_path = path.join(__dirname, "../.hfc-key-store");
console.log("Store path:" + store_path);
//const tx_id = null;

const userContextPromise = Fabric_Client.newDefaultKeyValueStore({path: store_path })
  .then(state_store => {

    console.log("Loading Credentials")
    fabric_client.setStateStore(state_store);

    const crypto_suite = Fabric_Client.newCryptoSuite();
    const crypto_store = Fabric_Client.newCryptoKeyStore({ path: store_path });

    crypto_suite.setCryptoKeyStore(crypto_store);
    fabric_client.setCryptoSuite(crypto_suite);

    return fabric_client.getUserContext("user1", true);
  });

async function query(serial) {
  
  return userContextPromise.then(async function(userContext) {
    
    if (userContext && userContext.isEnrolled()) {
      member_user = userContext;
    } else {
      throw new Error("Failed to get user1.... run registerUser.js");
    }

    const request = {
      chaincodeId: "misterybox",
      fcn:  "queryAllMisteryboxes",
      args: ""
    };
  
      // send the query proposal to the peer
      let buffer = channel.queryByChaincode(request);
      return await buffer.then(response => { return response[0].toString(); }).catch(err => {console.log(err)});
  }, error => {console.log("Erro ao obter credenciais" + error)});
}

async function insert(arr) {

  return userContextPromise.then(async function(userContext) {

    if (userContext && userContext.isEnrolled()) {
      member_user = userContext;
    } else {
      throw new Error("Failed to get user1.... run registerUser.js");
    }
  
    let tx_id = ""
    tx_id = await fabric_client.newTransactionID();

  
    const request = {
      chaincodeId: "misterybox",
      fcn: "createMisterybox",
      args: arr,
      chainId: "jewelchannel",
      txId: tx_id
    };
  
    let results = await channel.sendTransactionProposal(request);
    const proposalResponses = results[0];
    const proposal = results[1];
    let isProposalGood = false;
  
    if (proposalResponses && proposalResponses[0].response && proposalResponses[0].response.status === 200) {
      isProposalGood = true;
    } else {
      console.error("Transaction proposal was bad - " + proposalResponses);
    }
  
    if (isProposalGood) {
      console.log(
         util.format(
           'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s"',
           proposalResponses[0].response.status, proposalResponses[0].response.message ));
  
      const request = {
        proposalResponses: proposalResponses,
        proposal: proposal
      };

      const transaction_id_string = tx_id.getTransactionID();
      const promises = [];
  
      const sendPromise = channel.sendTransaction(request);
      promises.push(sendPromise); //we want the send transaction first, so that we know where to check status
  
      let event_hub = channel.newChannelEventHub(peer);
  
      let txPromise = new Promise((resolve, reject) => {
        
        let handle = setTimeout(() => {
          event_hub.unregisterTxEvent(transaction_id_string);
          event_hub.disconnect();
          resolve({ event_status: "TIMEOUT" });
        }, 3000);

        event_hub.registerTxEvent(
          transaction_id_string,
          (tx, code, blocknumber) => {
            clearTimeout(handle);


            console.log(`${blocknumber} :: ${code} :: ${tx} `)

  
            // now let the application know what happened
            const return_status = {
              event_status: code,
              tx_id: transaction_id_string,
              block_number: blocknumber,
              shim: proposalResponses[0].response.payload
            };
            if (code !== "VALID") {
              console.error("The transaction was invalid, code = " + code);
              resolve(return_status); // we could use reject(new Error('Problem with the tranaction, event status ::'+code));
            } else {
              //console.log("The transaction has been committed on peer " + event_hub.getPeerAddr());
              resolve(return_status);
            }
          },
          err => {
            //this is the callback if something goes wrong with the event registration or processing
            reject(new Error("There was a problem with the eventhub ::" + err));
          },
          { disconnect: true } //disconnect when complete
        );

        event_hub.connect();
      });
  
      promises.push(txPromise);
      return Promise.all(promises);
    }
    return proposalResponses[0].message; //Error
  });
  //const Transaction = () => {
  //    this._transaction_id = tx_id._transaction_id;
  //};
  //return Transaction;
}

async function update(arr) {
  
  return userContextPromise.then(async function(userContext) {

    if (userContext && userContext.isEnrolled()) {
      member_user = userContext;
    } else {
      throw new Error("Failed to get user1.... run registerUser.js");
    }

    let tx_id = ""
    tx_id = await fabric_client.newTransactionID();
    console.log("Assigning transaction_id: ", tx_id._transaction_id);

    const request = {

      chaincodeId: "misterybox",
      fcn: "transferMisterybox",
      args: arr,
      chainId: "jewelchannel",
      txId: tx_id
    };

    let results = await channel.sendTransactionProposal(request);
    const proposalResponses = results[0];
    const proposal = results[1];
    let isProposalGood = false;

    if (
      proposalResponses &&
      proposalResponses[0].response &&
      proposalResponses[0].response.status === 200
    ) {
      isProposalGood = true;
      console.log("Transaction proposal was good");
    } else{
      console.error("Transaction proposal was bad. " + proposalResponses + " --> " + arr);
    }

    if (isProposalGood) {
      console.log(
        util.format(
          'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s"',
          proposalResponses[0].response.status,
          proposalResponses[0].response.message
        )
      );

      const request = {
        proposalResponses: proposalResponses,
        proposal: proposal
      };
      const transaction_id_string = tx_id.getTransactionID();
      const promises = [];

      const sendPromise = channel.sendTransaction(request);
      promises.push(sendPromise); //we want the send transaction first, so that we know where to check status

      let event_hub = channel.newChannelEventHub(peer);

      let txPromise = new Promise((resolve, reject) => {
        let handle = setTimeout(() => {
          event_hub.unregisterTxEvent(transaction_id_string);
          event_hub.disconnect();
          resolve({ event_status: "TIMEOUT" });
        }, 3000);
        event_hub.registerTxEvent(
          transaction_id_string,
          (tx, code) => {
            // this is the callback for transaction event status
            // first some clean up of event listener
            clearTimeout(handle);

            // now let the application know what happened
            const return_status = {
              event_status: code,
              tx_id: transaction_id_string,
              shim: proposalResponses[0].response.payload
            };
            if (code !== "VALID") {
              console.error("The transaction was invalid, code = " + code);
              resolve(return_status); // we could use reject(new Error('Problem with the tranaction, event status ::'+code));
            } else {
              //console.log("The transaction has been committed on peer " + event_hub.getPeerAddr());
              resolve(return_status);
            }
          },
          err => {
            //this is the callback if something goes wrong with the event registration or processing
            reject(new Error("There was a problem with the eventhub ::" + err));
          },
          { disconnect: true } //disconnect when complete
        );
        event_hub.connect();
      });

      promises.push(txPromise);
      return Promise.all(promises);
    }
    return proposalResponses[0].message; //Error
  });
  

}

async function queryRestrictions(serial) {
  
  return userContextPromise.then(async function(userContext) {

    if (userContext && userContext.isEnrolled()) {
      member_user = userContext;
    } else {
      throw new Error("Failed to get user1.... run registerUser.js");
    }

    const request = {
      chaincodeId: "setupbox",
      fcn: serial != undefined ? "queryRestriction" : "queryAllRestrictions",
      args: [serial != undefined ? serial : ""]
    };

    // send the query proposal to the peer
    let buffer = channel.queryByChaincode(request);
    return await buffer.then(response => { return response[0].toString(); });
  });
}

module.exports.Query = query;
module.exports.New = insert;
module.exports.Update = update;
