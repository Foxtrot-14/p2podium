import {createRoot} from 'react-dom/client'
import './style.css'
import App from './App'
import { HashRouter } from 'react-router-dom'
import { Suspense } from 'react'

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <HashRouter>
      <Suspense>
        <App/>
      </Suspense>
    </HashRouter>
)
