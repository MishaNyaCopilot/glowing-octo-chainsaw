
import React, { useState, useEffect } from 'react';
import axios from 'axios';

const EpisodeForm = ({ episode, onSave, animeId, onCancel }) => {
  const [title, setTitle] = useState('');
  const [episodeNumber, setEpisodeNumber] = useState(1);

  useEffect(() => {
    if (episode) {
      setTitle(episode.title);
      setEpisodeNumber(episode.episode_number);
    } else {
      setTitle('');
      setEpisodeNumber(1);
    }
  }, [episode]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const episodeData = { title, episode_number: parseInt(episodeNumber, 10), anime_id: animeId };

    try {
      if (episode) {
        await axios.put(`http://localhost:8080/api/admin/episode/${episode.id}`, episodeData);
      } else {
        await axios.post('http://localhost:8080/api/admin/episode', episodeData);
      }
      onSave();
    } catch (error) {
      console.error('Error saving episode:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>Title</label>
        <input type="text" value={title} onChange={(e) => setTitle(e.target.value)} />
      </div>
      <div>
        <label>Episode Number</label>
        <input type="number" value={episodeNumber} onChange={(e) => setEpisodeNumber(e.target.value)} />
      </div>
      <button type="submit">Save</button>
      <button type="button" onClick={onCancel}>Cancel</button>
    </form>
  );
};

export default EpisodeForm;
