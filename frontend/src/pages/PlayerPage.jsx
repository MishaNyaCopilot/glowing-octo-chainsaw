
import React, { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import axios from 'axios';
import VideoPlayer from '../components/VideoPlayer';
import MainLayout from '../components/MainLayout';

const PlayerPage = () => {
  const { episodeId } = useParams();
  const [episode, setEpisode] = useState(null);
  const [anime, setAnime] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchEpisodeAndAnime = async () => {
      try {
        const episodeResponse = await axios.get(`http://localhost:8080/api/episode/${episodeId}`);
        setEpisode(episodeResponse.data);

        const animeResponse = await axios.get(`http://localhost:8080/api/anime/${episodeResponse.data.anime_id}`);
        if (animeResponse.data.episodes) {
          animeResponse.data.episodes.sort((a, b) => a.episode_number - b.episode_number);
        }
        setAnime(animeResponse.data);
      } catch (error) {
        console.error('Error fetching episode or anime details:', error);
      }
    };

    fetchEpisodeAndAnime();
  }, [episodeId]);

  if (!episode || !anime) {
    return <div>Loading...</div>;
  }

  const currentEpisodeIndex = anime.episodes.findIndex((e) => e.id === episode.id);
  const nextEpisode = anime.episodes[currentEpisodeIndex + 1];
  const prevEpisode = anime.episodes[currentEpisodeIndex - 1];

  const hlsUrl = `http://localhost:8080/hls/${episodeId}/playlist.m3u8`;

  return (
    <MainLayout title={episode.title}>
      <div className="w-full h-full">
        <VideoPlayer src={hlsUrl} />
      </div>
      <div className="flex justify-between mt-4">
        {prevEpisode && (
          <Link to={`/player/${prevEpisode.id}`}>
            <button className="bg-gray-800 text-white px-4 py-2 rounded-lg">Previous Episode</button>
          </Link>
        )}
        {nextEpisode && (
          <Link to={`/player/${nextEpisode.id}`}>
            <button className="bg-gray-800 text-white px-4 py-2 rounded-lg">Next Episode</button>
          </Link>
        )}
      </div>
      <div className="mt-8">
        <h2 className="text-2xl font-bold mb-4">Episodes</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {anime.episodes.map((ep) => (
            <Link to={`/player/${ep.id}`} key={ep.id}>
              <div className={`p-4 rounded-lg ${ep.id === episode.id ? 'bg-gray-700' : 'bg-gray-800 hover:bg-gray-700'}`}>
                Episode {ep.episode_number}: {ep.title}
              </div>
            </Link>
          ))}
        </div>
      </div>
    </MainLayout>
  );
};

export default PlayerPage;
