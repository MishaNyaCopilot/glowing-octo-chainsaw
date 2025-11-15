
import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import axios from 'axios';
import AnimeForm from '../components/admin/AnimeForm';
import EpisodeForm from '../components/admin/EpisodeForm';
import UploadForm from '../components/admin/UploadForm';
import UploadPosterForm from '../components/admin/UploadPosterForm';

const AdminPage = () => {
  const [animes, setAnimes] = useState([]);
  const [filteredAnimes, setFilteredAnimes] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedAnime, setSelectedAnime] = useState(null);
  const [selectedEpisode, setSelectedEpisode] = useState(null);
  const [showAnimeForm, setShowAnimeForm] = useState(false);
  const [showEpisodeForm, setShowEpisodeForm] = useState(false);
  const [showUploadForm, setShowUploadForm] = useState(false);
  const [showUploadPosterForm, setShowUploadPosterForm] = useState(false);

  useEffect(() => {
    fetchAnimes();
  }, []);

  useEffect(() => {
    setFilteredAnimes(
      animes.filter((anime) =>
        anime.title.toLowerCase().includes(searchTerm.toLowerCase())
      )
    );
  }, [searchTerm, animes]);

  const fetchAnimes = async () => {
    try {
      const response = await axios.get('http://localhost:8080/api/anime/');
      setAnimes(response.data);
    } catch (error) {
      console.error('Error fetching animes:', error);
    }
  };

  const fetchAnime = async (id) => {
    try {
      const response = await axios.get(`http://localhost:8080/api/anime/${id}`);
      // Sort episodes by episode_number
      if (response.data.episodes) {
        response.data.episodes.sort((a, b) => a.episode_number - b.episode_number);
      }
      setSelectedAnime(response.data);
      return response.data;
    } catch (error) {
      console.error('Error fetching anime details:', error);
    }
  };

  const handleSaveAnime = () => {
    setShowAnimeForm(false);
    setSelectedAnime(null);
    fetchAnimes();
  };

  const handleSaveEpisode = () => {
    setShowEpisodeForm(false);
    setSelectedEpisode(null);
    if (selectedAnime) {
      fetchAnime(selectedAnime.id);
    }
  };

  const handleUpload = () => {
    setShowUploadForm(false);
    setSelectedEpisode(null);
  };

  const handleUploadPoster = () => {
    setShowUploadPosterForm(false);
    setSelectedAnime(null);
    fetchAnimes();
  };

  const handleDeleteAnime = async (id) => {
    try {
      await axios.delete(`http://localhost:8080/api/admin/anime/${id}`);
      fetchAnimes();
      setSelectedAnime(null);
    } catch (error) {
      console.error('Error deleting anime:', error);
    }
  };

  const handleDeleteEpisode = async (id) => {
    try {
      await axios.delete(`http://localhost:8080/api/admin/episode/${id}`);
      if (selectedAnime) {
        fetchAnime(selectedAnime.id);
      }
    } catch (error) {
      console.error('Error deleting episode:', error);
    }
  };

  const handleEditClick = async (anime) => {
    const fetchedAnime = await fetchAnime(anime.id);
    setSelectedAnime(fetchedAnime);
    setShowAnimeForm(true);
  };

  const handleUploadPosterClick = async (anime) => {
    const fetchedAnime = await fetchAnime(anime.id);
    setSelectedAnime(fetchedAnime);
    setShowUploadPosterForm(true);
  };

  return (
    <div className="container">
      <Link to="/" className="nav-button">Go to Main Page</Link>
      <h1>Admin Panel</h1>

      {showAnimeForm && (
        <div className="modal">
          <div className="modal-content">
            <span className="close" onClick={() => setShowAnimeForm(false)}>&times;</span>
            <AnimeForm anime={selectedAnime} onSave={handleSaveAnime} onCancel={() => setShowAnimeForm(false)} />
          </div>
        </div>
      )}

      {showEpisodeForm && (
        <div className="modal">
          <div className="modal-content">
            <span className="close" onClick={() => setShowEpisodeForm(false)}>&times;</span>
            <EpisodeForm episode={selectedEpisode} onSave={handleSaveEpisode} animeId={selectedAnime.id} onCancel={() => setShowEpisodeForm(false)} />
          </div>
        </div>
      )}

      {showUploadForm && (
        <div className="modal">
          <div className="modal-content">
            <span className="close" onClick={() => setShowUploadForm(false)}>&times;</span>
            <UploadForm episodeId={selectedEpisode.id} onUpload={handleUpload} />
          </div>
        </div>
      )}

      {showUploadPosterForm && selectedAnime && (
        <div className="modal">
          <div className="modal-content">
            <span className="close" onClick={() => setShowUploadPosterForm(false)}>&times;</span>
            <UploadPosterForm animeId={selectedAnime.id} onUpload={handleUploadPoster} onUploadSuccess={fetchAnimes} />
          </div>
        </div>
      )}

      <h2>Anime Management</h2>
      <div className="toolbar">
        <input
          type="text"
          placeholder="Search for an anime..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="search-input"
        />
        <button onClick={() => { setSelectedAnime(null); setShowAnimeForm(true); }}>Add Anime</button>
      </div>
      <table>
        <thead>
          <tr>
            <th>Title</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {filteredAnimes.map((anime) => (
            <tr key={anime.id}>
              <td onClick={() => fetchAnime(anime.id)} style={{ cursor: 'pointer' }}>
                {anime.title}
              </td>
              <td>
                <button onClick={() => handleEditClick(anime)}>Edit</button>
                <button onClick={() => handleDeleteAnime(anime.id)}>Delete</button>
                <button onClick={() => handleUploadPosterClick(anime)}>Upload Poster</button>
                <button onClick={() => fetchAnime(anime.id)}>Manage Episodes</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {selectedAnime && (
        <div>
          <h2>Episodes for {selectedAnime.title}</h2>
          <button onClick={() => { setSelectedEpisode(null); setShowEpisodeForm(true); }}>Add Episode</button>
          <table>
            <thead>
              <tr>
                <th>Episode Number</th>
                <th>Title</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {selectedAnime.episodes && selectedAnime.episodes.map((episode) => (
                <tr key={episode.id}>
                  <td>{episode.episode_number}</td>
                  <td>{episode.title}</td>
                  <td>
                    <button onClick={() => { setSelectedEpisode(episode); setShowEpisodeForm(true); }}>Edit</button>
                    <button onClick={() => handleDeleteEpisode(episode.id)}>Delete</button>
                    <button onClick={() => { setSelectedEpisode(episode); setShowUploadForm(true); }}>Upload Video</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default AdminPage;
