import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import HomePage from './pages/HomePage';
import AnimeDetailPage from './pages/AnimeDetailPage';
import PlayerPage from './pages/PlayerPage';
import AdminPage from './pages/AdminPage'; // Import the AdminPage component
import './App.css';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/anime/:id" element={<AnimeDetailPage />} />
        <Route path="/player/:episodeId" element={<PlayerPage />} />
        <Route path="/admin" element={<AdminPage />} /> {/* Add the admin route */}
      </Routes>
    </Router>
  );
}

export default App;
