import { ImageResponse } from 'next/server'
// App router includes @vercel/og.
// No need to install it.
import base from '../../public/base.png'
export const runtime = 'edge';

export async function GET() {
  console.log(base)
    return new ImageResponse(
      (
        <div
          style={{
            width: '100%',
            height: '100%',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontSize: 128,
            background: 'lavender',
            // backgroundImage: `url(${}),`
          }}
        >
          
        </div>
      )
    )
  }