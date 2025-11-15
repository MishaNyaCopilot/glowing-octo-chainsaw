
import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import axios from 'axios';
import MainLayout from '../components/MainLayout';
import EpisodeCard from '../components/EpisodeCard';

const AnimeDetailPage = () => {
  const { id } = useParams();
  const [anime, setAnime] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchAnime = async () => {
      try {
        const response = await axios.get(`http://localhost:8080/api/anime/${id}`);
        // Sort episodes by episode_number
        if (response.data.episodes) {
          response.data.episodes.sort((a, b) => a.episode_number - b.episode_number);
        }
        setAnime(response.data);
      } catch (error) {
        setError('Error fetching anime details. Please try again later.');
        console.error('Error fetching anime details:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchAnime();
  }, [id]);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>{error}</div>;
  }

  return (
    <MainLayout title={anime.title}>
      <div className="flex flex-col md:flex-row gap-8">
        <div className="md:w-1/3">
          <img src={`http://localhost:8080${anime.poster_url}` || 'https://via.placeholder.com/300x450'} alt={anime.title} className="h-auto rounded-lg" />
        </div>
        <div className="md:w-2/3">
          <p className="text-lg">{anime.description}</p>
          <h2 className="text-2xl font-bold mt-8 mb-4">Episodes</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {anime.episodes && anime.episodes.length > 0 ? (
              anime.episodes.map((episode) => <EpisodeCard key={episode.id} episode={episode} />)
            ) : (
              <div>No episodes found.</div>
            )}
          </div>
        </div>
      </div>
    </MainLayout>
  );
};

export default AnimeDetailPage;
