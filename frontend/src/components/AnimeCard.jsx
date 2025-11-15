import React from 'react';
import { Link } from 'react-router-dom';

const AnimeCard = ({ anime }) => {
  return (
    <Link to={`/anime/${anime.id}`}>
      <div className="bg-gray-800 rounded-lg overflow-hidden">
        <img src={`http://localhost:8080${anime.poster_url}` || 'https://via.placeholder.com/300x450'} alt={anime.title} className="h-96 object-cover object-center" />
        <div className="p-4">
          <h3 className="text-lg font-bold">{anime.title}</h3>
          <p className="text-sm text-gray-400 truncate">{anime.description}</p>
        </div>
      </div>
    </Link>
  );
};

export default AnimeCard;