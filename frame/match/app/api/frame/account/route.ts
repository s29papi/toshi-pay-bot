import { FrameRequest, getFrameMessage } from '@coinbase/onchainkit';
import { NextRequest, NextResponse } from 'next/server';

async function getResponse(req: NextRequest): Promise<NextResponse> {
  const body: FrameRequest = await req.json();
  const searchParams = req.nextUrl.searchParams;
  const gameId:any = searchParams.get("gameId");
  const buttonId = body.untrustedData.buttonIndex;
  // pass in facaster Id + hash of transaction generated from the tx
  // bot would resolve the balance
  const { isValid, message } = await getFrameMessage(body , {
    neynarApiKey: 'NEYNAR_ONCHAIN_KIT'
  });
   // back 
  if (buttonId == 1) {
    return NextResponse.redirect("https://wag3r-bot.vercel.app/?gameId=${gameId}");
  }
  // withdraw
  if (buttonId == 2) {
    return NextResponse.redirect('https://wag3r-bot.vercel.app/~/unstake?${message.fid}', {status: 302});
  }
  // deposit
  if (buttonId == 3) {
    return NextResponse.redirect('https://wag3r-bot.vercel.app/~/stake?${message.fid}', {status: 302});
  }

  // refresh 
  return new NextResponse(`<!DOCTYPE html><html><head>
        <title>Account</title>
        <meta property="fc:frame" content="vNext" />
        <meta property="fc:frame:image" content="https://wag3r-bot.vercel.app/og/account"/>
        <meta property="fc:frame:button:1" content="Back" />
        <meta property="fc:frame:button:1:action" content="post"/>
        <meta property="fc:frame:button:2" content="Withdraw" />
        <meta property="fc:frame:button:2:action" content="post_redirect"/>
        <meta property="fc:frame:button:3" content="Deposit" />
        <meta property="fc:frame:button:3:action" content="post_redirect"/>
        <meta property="fc:frame:button:4" content="Refresh" />
        <meta property="fc:frame:button:4:action" content="post_redirect"/>
        <meta property="fc:frame:post_url" content="https://wag3r-bot.vercel.app/api/frame/account?gameId=${gameId}"/>
        </head></html>`);
}

export async function POST(req: NextRequest): Promise<Response> {
  return getResponse(req);
}

export const dynamic = 'force-dynamic';

