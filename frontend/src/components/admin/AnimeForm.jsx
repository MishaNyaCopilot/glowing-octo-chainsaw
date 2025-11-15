
import React, { useState, useEffect } from 'react';
import axios from 'axios';

const AnimeForm = ({ anime, onSave, onCancel }) => {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');

  useEffect(() => {
    if (anime) {
      setTitle(anime.title);
      setDescription(anime.description);
    } else {
      setTitle('');
      setDescription('');
    }
  }, [anime]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const animeData = { title, description };

    try {
      if (anime) {
        await axios.put(`http://localhost:8080/api/admin/anime/${anime.id}`, animeData);
      } else {
        await axios.post('http://localhost:8080/api/admin/anime', animeData);
      }
      onSave();
    } catch (error) {
      console.error('Error saving anime:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>Title</label>
        <input type="text" value={title} onChange={(e) => setTitle(e.target.value)} />
      </div>
      <div>
        <label>Description</label>
        <textarea value={description} onChange={(e) => setDescription(e.target.value)} />
      </div>
      <button type="submit">Save</button>
      <button type="button" onClick={onCancel}>Cancel</button>
    </form>
  );
};

export default AnimeForm;
