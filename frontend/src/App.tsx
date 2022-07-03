import React from 'react';
import Coin from "./components/Coin"
import NFT from "./components/NFT"

import './App.css';

function App() {
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
      <div>
        <NFT
          signer_address="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
          friend_address="0x70997970c51812dc3a010c7d01b50e0d17dc79c8"
        />
      </div>
      <div>
        <NFT
          signer_address="0x70997970c51812dc3a010c7d01b50e0d17dc79c8"
          friend_address="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
        />
      </div>
    </>
  )
}

export default App;
