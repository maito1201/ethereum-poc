import React from 'react';
import Poll from "./components/Poll"
import './App.css';
import { ethers } from "ethers";
import artifact from "./abi/Poll.json";

function App() {
  const contractAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3"
  const provider = new ethers.providers.JsonRpcProvider();
  const signer = provider.getSigner();
  const contract = new ethers.Contract(contractAddress, artifact.abi, provider);
  const contractWithSigner = contract.connect(signer);

  return (
    <div>
      <h1>Hello, Poll Contract.</h1>
      <Poll contract={contractWithSigner} />
    </div>
  )
}

export default App;
