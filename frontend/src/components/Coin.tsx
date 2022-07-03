import React, { useState, useEffect } from 'react';
import { ethers } from "ethers";
import artifact from "../abi/MyCoin.json"

function Coin(props: { signer_address: string, friend_address: string }) {
  const provider = new ethers.providers.JsonRpcProvider();
  const coin_address = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
  const signer = provider.getSigner(props.signer_address);
  const contract = new ethers.Contract(coin_address, artifact.abi, provider);
  const client = contract.connect(signer);

  const [isOwner, setIsOwner] = useState<boolean>(false)
  const [balance, setBalance] = useState<string>("")
  const [coinName, setCoinName] = useState<string>("")
  const [symbolName, setSymbolName] = useState<string>("")
  const { balanceOf, name, symbol, mint, approve, transfer, transferFrom } = client.functions;
  useEffect(() => {
    const setNameFunc = async () => {
      setCoinName((await name()).toString())
      setSymbolName((await symbol()).toString())
    }
    const setIsOwnerFunc = async () => {
      const owner = (await client.owner()).toString()
      const signer = (await client.signer.getAddress()).toString()
      setIsOwner(owner === signer)
    } 
    setNameFunc()
    setIsOwnerFunc()
    setBalanceFunc()
  }, [])

  const setBalanceFunc = async () => {
    const ret = (await balanceOf(props.signer_address)).toString()
    setBalance(await ret)
  }
  
  const hundleMint = async () => {
    mint(props.signer_address, 100).then(() => {
      setBalanceFunc()
    }).catch((error: any) => {
      alert(error)
    })
  }

  const hundleApprove = async() => {
    await approve(props.friend_address, 100)
  }

  const hundleTransferFromFriend = async () => {
    transferFrom(props.friend_address, props.signer_address, 100).then(() =>{
      setBalanceFunc()
    }).catch((error: any) => {
      alert(error)
    })
  }

  const hundleTransfer = async() => {
    transfer(props.friend_address, 100).then(() => {
      setBalanceFunc()
    }).catch((error: any) => {
      alert(error)
    })
  }

  return (
    <>
    <div style={{margin: "40px"}}>
      <div>{coinName}</div>
      <div>user_address: {props.signer_address}{isOwner ? " (owner)" : ""}</div>
      <div>amount: {balance}{symbolName}</div>
      <div>
        <button onClick={() => { hundleTransfer() }}>transfer to friend</button>
        <button onClick={() => { hundleApprove() }}>approve to transfer from friend</button>
        <button onClick={() => { hundleTransferFromFriend() }}>transfer from friend</button>
      </div>
      {isOwner ?
        <div>
          <button onClick={() => { hundleMint() }}>mint</button>
        </div> :
        null
      }
    </div>
    </>
  )
}

export default Coin;
