import React from 'react';
import Coin from "./components/Coin"

import './App.css';
// import nft_artifact from "./abi/MyNFT.json"

function App() {
  const poll_address = "0x5FbDB2315678afecb367f032d93F642f64180aa3"
  // const nft_address = "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0"

  return (
    <>
      <div>
        <Coin
          signer_address="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
          friend_address="0x70997970c51812dc3a010c7d01b50e0d17dc79c8"
        />
      </div>
      <div>
        <Coin
          signer_address="0x70997970c51812dc3a010c7d01b50e0d17dc79c8"
          friend_address="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
        />
      </div>
    </>
  )
}

export default App;
