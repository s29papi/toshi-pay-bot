'use client';

import {useRouter} from "next/navigation";
import {useEffect} from "react";
import Connect from '../../components/Connect'
// import { useWeb3Modal } from '@web3modal/wagmi/react'
import { useSendTransaction, useSignMessage  } from 'wagmi' 
import { parseEther } from 'viem' 



export default function Redirect() {
    const router = useRouter();
    // const { open } = useWeb3Modal()
    const { data: hash, sendTransaction } = useSendTransaction() 
    const { signMessage } = useSignMessage()

   

    useEffect(() => {
        const handleBeforeUnload = (event: BeforeUnloadEvent) => {
          // Cancel the default behavior of closing the tab
          event.preventDefault();
          // Chrome requires the returnValue to be set
          event.returnValue = '';
        };
    
        // Add event listener to beforeunload event
        window.addEventListener('beforeunload', handleBeforeUnload);
    
        return () => {
          // Remove the event listener when component unmounts
          window.removeEventListener('beforeunload', handleBeforeUnload);
        };
      }, []);

    // useEffect(() => {
    //     async function openWallet() {
    //       // open()
          
  
    //     }

    //     openWallet()
    //     sendTransaction({to: `0x${"47dEAF612F0769d99aDB653bA2d22bba79F26C42"}`, value: parseEther("0.1")})
    // })

    function submitTx() {
      sendTransaction({to: `0x${"47dEAF612F0769d99aDB653bA2d22bba79F26C42"}`, value: parseEther("0.1")})
    }

      const handleCloseButtonClick = () => {
        // Close the current tab
        window.close();
      };


    return (
        <div>
            
            <Connect />
            <button onClick={() => { submitTx() }}>Stake</button>
            <button onClick={handleCloseButtonClick}>Close Tab</button>
            <button onClick={() => signMessage({ message: 'hello world' })}>Sign message</button>
            <button onClick={() => sendTransaction({ to: `0x${"47dEAF612F0769d99aDB653bA2d22bba79F26C42"}`, value: parseEther("0.2") })}>Send Tx</button>
        </div>
    );
}



// exit page when done