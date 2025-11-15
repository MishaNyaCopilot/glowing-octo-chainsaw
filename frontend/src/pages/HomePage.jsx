import React, { useState, useEffect } from 'react';
import axios from 'axios';
import MainLayout from '../components/MainLayout';
import AnimeCard from '../components/AnimeCard';

const HomePage = () => {
  const [animes, setAnimes] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchAnimes = async () => {
      try {
        const response = await axios.get('http://localhost:8080/api/anime/');
        setAnimes(response.data);
      } catch (error) {
        setError('Error fetching animes. Please try again later.');
        console.error('Error fetching animes:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchAnimes();
  }, []);

  const filteredAnimes = animes.filter((anime) =>
    anime.title.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <MainLayout>
      <h1 className="text-4xl font-bold text-center mb-8">AniStream</h1>
      <div className="mb-4">
        <input
          type="text"
          placeholder="Search for an anime..."
          className="w-full p-2 rounded bg-gray-800 text-white"
          onChange={(e) => setSearchTerm(e.target.value)}
        />
      </div>
      {loading && <div>Loading...</div>}
      {error && <div>{error}</div>}
      {!loading && !error && (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {filteredAnimes.length > 0 ? (
            filteredAnimes.map((anime) => <AnimeCard key={anime.id} anime={anime} />)
          ) : (
            <div className="col-span-full text-center">No results found.</div>
          )}
        </div>
      )}
    </MainLayout>
  );
};

export default HomePage;
