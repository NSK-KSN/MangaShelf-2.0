import React, { useState, useEffect } from 'react';
import axios from 'axios';

function CoverImageUploader() {
  const [mangaTitles, setMangaTitles] = useState([]);
  const [selectedTitle, setSelectedTitle] = useState('');
  const [selectedFile, setSelectedFile] = useState(null);
  const [volumeNumber, setVolumeNumber] = useState('');
  
  useEffect(() => {
    // Fetch manga titles from the backend
    const fetchMangaTitles = async () => {
      try {
        const response = await axios.get('http://localhost:8080/fetch-manga-titles');
        setMangaTitles(response.data);
        console.log(response.data);
      } catch (error) {
        console.error('Error fetching manga titles:', error);
      }
    };

    fetchMangaTitles();
  }, []);

  const handleFileChange = (event) => {
    setSelectedFile(event.target.files[0]);
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    
    if (!selectedTitle || !selectedFile || !volumeNumber) {
      console.error('Please select a title, file, and enter a volume number.');
      return;
    }

    const formData = new FormData();
    formData.append('file', selectedFile);
    formData.append('title', selectedTitle);
    formData.append('volumeNumber', volumeNumber);

    try {
      await axios.post('http://localhost:8080/upload-cover-image', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      });
      console.log('Cover image uploaded successfully.');
    } catch (error) {
      console.error('Error uploading cover image:', error);
    }
  };

  return (
    <div>
      <h1>Upload a cover image</h1>
    <form onSubmit={handleSubmit}>
      <select value={selectedTitle} onChange={(e) => setSelectedTitle(e.target.value)}>
        <option value="">Select Manga Title</option>
        {mangaTitles.map((title, index) => (
          <option key={index} value={title}>{title}</option>
        ))}
      </select>
      <input type="file" onChange={handleFileChange} />
      <input
        type="text"
        value={volumeNumber}
        onChange={(e) => setVolumeNumber(e.target.value)}
        placeholder="Volume Number"
      />
      <button type="submit">Upload Cover Image</button>
    </form>
    </div>
  );
}

export default CoverImageUploader;