import Versions from './components/Versions'
import electronLogo from './assets/electron.svg'
import { isElectron } from './utils/electron'

function App(): JSX.Element {
  const ipcHandle = (): void => window.electron.ipcRenderer.send('ping')

  const test = (): void => {
    fetch('http://127.0.0.1:3000/api/users').then((res) => res.text())
  }

  return (
    <div
      className="flex flex-col items-center justify-center"
      style={{ backgroundColor: 'var(--color-background)' }}
    >
      <div
        className="p-6 text-center w-full md:w-1/2 lg:w-1/3"
        style={{ backgroundColor: 'var(--ev-c-gray-2)', color: 'var(--ev-c-white)' }}
      >
        <h1 className="text-2xl md:text-xl lg:text-3xl">Hello Tailwind!</h1>
      </div>
      <img alt="logo" className="w-32 h-32 my-6" src={electronLogo} />
      <div className="text-base md:text-sm lg:text-base" style={{ color: 'var(--ev-c-text-1)' }}>
        Powered by electron-vite
      </div>
      <div
        className="text-lg md:text-base lg:text-lg text-center my-6"
        style={{ color: 'var(--ev-c-text-2)' }}
      >
        Build an Electron app with <span className="font-semibold">React</span> &nbsp; and{' '}
        <span className="font-semibold">TypeScript</span>
      </div>
      <p className="text-sm md:text-xs lg:text-sm" style={{ color: 'var(--ev-c-text-3)' }}>
        Please try pressing <code>F12</code> to open the devTool
      </p>
      <div className="flex flex-col items-center space-y-4 mt-6">
        <div>
          <a
            href="https://electron-vite.org/"
            target="_blank"
            rel="noreferrer"
            className="hover:underline text-lg"
            style={{ color: 'var(--ev-c-text-1)' }}
          >
            Documentation
          </a>
        </div>
        <div>
          <button
            onClick={ipcHandle}
            className="focus:outline-none hover:underline text-lg"
            style={{ color: 'var(--ev-c-text-1)' }}
          >
            Send IPC
          </button>
        </div>
        <div>
          <button
            onClick={test}
            className="focus:outline-none hover:underline text-lg"
            style={{ color: 'var(--ev-c-text-1)' }}
          >
            Test
          </button>
        </div>
      </div>
      {isElectron() && <Versions />}
    </div>
  )
}

export default App
