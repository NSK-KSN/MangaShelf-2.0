import Navbar from './components/navbar/Navbar.jsx';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Home from './components/pages/Home.jsx';
import Browse from './components/pages/Browse.jsx';
import MangaPage from './components/pages/MangaPage.jsx';
import AdminPage from './components/pages/AdminPage.jsx';
import './styles/App.css';
import React, { useEffect, useState } from 'react'

export default function App() {

  const current_theme = localStorage.getItem('current_theme');
  const [theme, setTheme] = useState(current_theme? current_theme : 'light');

  useEffect(() => {
    localStorage.setItem('current_theme', theme);
  }, [theme])

  return (
    <div className={`App ${theme}`}>
      <BrowserRouter>
      <Navbar theme={theme} setTheme={setTheme}/>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="browse" element={<Browse />} />
          <Route path="/manga/:id" element={<MangaPage/>} />
          <Route path="admin" element={<AdminPage />} />
        </Routes>
      </BrowserRouter>
    </div> 
  );
}
