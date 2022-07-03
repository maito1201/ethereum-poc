import React, { useState, useEffect } from 'react';
import { ethers } from "ethers";
import artifact from "../abi/MyNFT.json"

function NFT(props: { signer_address: string, friend_address: string }) {
  type TokenItem = {
    tokenID: string
    tokenURI: string
  }
  const provider = new ethers.providers.JsonRpcProvider();
  const nft_address = "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0"
  const signer = provider.getSigner(props.signer_address);
  const contract = new ethers.Contract(nft_address, artifact.abi, provider);
  const client = contract.connect(signer);

  const [balance, setBalance] = useState<string>("")
  const [isOwner, setIsOwner] = useState<boolean>(false)
  const [tokenName, setTokenName] = useState<string>("")
  const [tokenItems, setTokenItems] = useState<Array<TokenItem>>([])
  const { balanceOf, safeMint, name, transferFrom, tokenURI } = client.functions;

  useEffect(() => {
    const setTokenNameFunc = async () => {
      setTokenName((await name()).toString())
    }
    const setIsOwnerFunc = async () => {
      const owner = (await client.owner()).toString()
      const signer = (await client.signer.getAddress()).toString()
      setIsOwner(owner === signer)
    }

    setTokenNameFunc()
    setIsOwnerFunc()
    setBalanceFunc()
    setTokenItemsFunc()
  }, [])

  const setBalanceFunc = async () => {
    const ret = (await balanceOf(props.signer_address)).toString()
    setBalance(await ret)
  }

  const setTokenItemsFunc = async () => {
    let items: Array<TokenItem> = [];
    const sentLogs = await client.queryFilter(
      client.filters.Transfer(props.signer_address, null),
    );
    const receivedLogs = await client.queryFilter(
      client.filters.Transfer(null, props.signer_address),
    );

    const logs = sentLogs.concat(receivedLogs).sort((a, b) =>
      a.blockNumber - b.blockNumber ||
      a.transactionIndex - b.transactionIndex,
    )

    const owned = new Set<string>()
    for (const log of logs) {
      if (log?.args?.to === (await client.signer.getAddress()).toString()) {
        owned.add(log?.args?.tokenId.toString())
      } else {
        owned.delete(log?.args?.tokenId.toString())
      }
    }

    for (const id of owned) {
      const uri = (await tokenURI(id)).toString()
      const item: TokenItem = {
        tokenID: id,
        tokenURI: uri
      }
      items.push(item)
    }
    setTokenItems(items)
  }
  
  const hundleMint = async () => {
    safeMint(props.signer_address, "https://1.bp.blogspot.com/-drtIDNvYq80/YVwDI7wBDCI/AAAAAAABfdM/tRz-NyBz1C4DfXPQGkQ9uwBKdnHyFrBsQCNcBGAsYHQ/s400/internet_nft_art.png").then(() => {
      setBalanceFunc()
      setTokenItemsFunc()
    }).catch((error: any) => {
      alert(error)
    })
  }

  const hundleApprove = async() => {
    // await approve(props.friend_address, 100)
  }

  const hundleTransferFromFriend = async () => {
    /*
    transferFrom(props.friend_address, props.signer_address, 100).then(() =>{
      setBalanceFunc()
    }).catch((error: any) => {
      alert(error)
    })
    */
  }

  const hundleTransfer = async(id: string) => {
    transferFrom(props.signer_address, props.friend_address, id).then(() => {
      setBalanceFunc()
      setTokenItemsFunc()
    }).catch((error: any) => {
      alert(error)
    })
  }

  const tokenItem = (item: TokenItem) => {
    return(
      <div key={item.tokenID} style={{ margin: "10px 0" }}>
        <div>id: {item.tokenID}</div>
        <img src={item.tokenURI} style={{width: "200px"}} />
        <div>
          <button onClick={() => { hundleTransfer(item.tokenID) }}>transfer to friend</button>
        </div>
      </div>
    )
  }

  return (
    <>
    <div style={{margin: "10px 20px", borderRadius: "10px", boxShadow: "2px 2px 10px #bbb", padding: "20px"}}>
      <div>{tokenName}</div>
      <div>user_address: {props.signer_address}{isOwner ? " (owner)" : ""}</div>
      {isOwner ?
        <div>
          <button onClick={() => { hundleMint() }}>mint</button>
        </div> :
        null
      }
      <div>nfts: {balance}items</div>
      {tokenItems.map(tokenItem)}
    </div>
    </>
  )
}

export default NFT;
