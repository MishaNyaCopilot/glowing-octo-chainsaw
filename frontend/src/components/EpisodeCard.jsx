import React from 'react';
import { Link } from 'react-router-dom';

const EpisodeCard = ({ episode }) => {
  return (
    <Link to={`/player/${episode.id}`}>
      <div className="bg-gray-800 rounded-lg overflow-hidden transform transition-transform duration-300 hover:scale-105">
        <div className="p-4">
          <h3 className="text-lg font-bold">Episode {episode.episode_number}</h3>
          <p className="text-sm text-gray-400 truncate">{episode.title}</p>
        </div>
      </div>
    </Link>
  );
};

export default EpisodeCard;