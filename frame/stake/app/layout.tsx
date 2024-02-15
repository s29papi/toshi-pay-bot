import './globals.css'

import { Web3ModalProvider } from '../context/Web3Modal'

export const metadata = {
  title: 'Web3Modal',
  description: 'Web3Modal Example'
}

export default function RootLayout({ children }: any) {
  return (
    <html lang="en">
      <body>
        <Web3ModalProvider>{children}</Web3ModalProvider>
      </body>
    </html>
  )
}