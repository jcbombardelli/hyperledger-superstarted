const express = require("express");
const bodyParser = require("body-parser");
const app = express();

const timeout = require('connect-timeout')

app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());

app.use(timeout('45s'))

const blockchain = require("../infraestructure/js/00-lib");

app.get("/api/misterybox", async (req, res) => {
  res.setHeader("Content-Type", "application/json");
  blockchain.Query().then(queryResult => {
    res.status(200).send(queryResult);
  });
});

app.post("/api/misterybox", async (req, res) => {
  res.setHeader("Content-Type", "application/json");
  let arr = [];
  arr[0] = req.body.serial;
  arr[1] = req.body.size;
  arr[2] = req.body.model;
  arr[3] = req.body.owner;
  
  blockchain.New(arr).then(result => {
    try {
      res.status(200).send(result);
    } catch (error) {
      console.log(`Blockchain Error: ${error} -[${arr}] result ${result}`);
      res.status(500).send(error);
    }
  });
});

app.put("/api/misterybox/:serial", async (req, res) => {
  res.setHeader("Content-Type", "application/json");

  let arr = [];
  arr[0] = req.param("serial");
  arr[1] = req.body.owner;
  arr[2] = req.body.newowner;
  
  blockchain.Update(arr).then(queryResult => {
    res.status(200).send(queryResult);
  })
})

app.listen(3000, function () {
  console.log("i3Tech * Hyperledger app listening on port 3000!");
});