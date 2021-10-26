import React, { Suspense } from 'react';

import { HashRouter } from 'react-router-dom'

import MainRouter from './Routes';

import './App.css';
import '@fontsource/roboto';

function TheApp() {

  return (
    <HashRouter>
      <MainRouter />
    </HashRouter>
  );
}


export default function App() {
  return (
    <Suspense fallback={<div>loading</div>}>
      <TheApp />
    </Suspense>
  )
};
