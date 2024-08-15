import Versions from './components/Versions'
import electronLogo from './assets/electron.svg'
import { isElectron } from './utils/electron'

function App(): JSX.Element {
  const ipcHandle = (): void => window.electron.ipcRenderer.send('ping')

  const test = (): void => {
    fetch('http://127.0.0.1:3000/api/users').then((res) => res.text())
  }

  return (
    <>
      <img alt="logo" className="logo" src={electronLogo} />
      <div className="creator">Powered by electron-vite</div>
      <div className="text">
        Build an Electron app with <span className="react">React</span>
        &nbsp;and <span className="ts">TypeScript</span>
      </div>
      <p className="tip">
        Please try pressing <code>F12</code> to open the devTool
      </p>
      <div className="actions">
        <div className="action">
          <a href="https://electron-vite.org/" target="_blank" rel="noreferrer">
            Documentation
          </a>
        </div>
        <div className="action">
          <a target="_blank" rel="noreferrer" onClick={ipcHandle}>
            Send IPC
          </a>
        </div>

        <div className="action">
          <a target="_blank" rel="noreferrer" onClick={test}>
            Test
          </a>
        </div>
      </div>

      {isElectron() ? <Versions /> : null}
    </>
  )
}

export default App
