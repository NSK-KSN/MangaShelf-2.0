import React, { useState, useEffect } from 'react';
import axios from 'axios';

function AddTitleForm() {
  const [formData, setFormData] = useState({
    title: '',
    cover_image: [], // Initialize as an empty array
    type_id: '',
    publisher_id: '',
    mal_id: '',
    score: '',
    popularity: '',
    status_id: ''
  });

  const [publishers, setPublishers] = useState([]);
  const [types, setTypes] = useState([]);
  const [status, setStatus] = useState([]);

  useEffect(() => {
    // Fetch publishers from backend
    const fetchPublishers = async () => {
      try {
        const publishersResponse = await axios.get('http://localhost:8080/dbquery/publishers');
        setPublishers(publishersResponse.data);

        const typesResponse = await axios.get('http://localhost:8080/dbquery/types');
        setTypes(typesResponse.data);

        const statusResponse = await axios.get('http://localhost:8080/dbquery/status');
        setStatus(statusResponse.data);
      } catch (error) {
        console.error('Error fetching publishers:', error);
      }
    };

    fetchPublishers();
  }, []);

  const [selectedFile, setSelectedFile] = useState(null);

  const handleFileChange = (e) => {
    setSelectedFile(e.target.files[0]);
  };

  const handleUploadAndAddData = async () => {
    // Upload the file
    const formDataFile = new FormData();
    formDataFile.append('file', selectedFile);
    try {
      const response = await axios.post('http://localhost:8080/upload', formDataFile, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      });
      console.log(response.data);
      // Handle successful upload
    } catch (error) {
      console.error('Error:', error);
      // Handle upload error
    }

    const fileName = selectedFile.name;
    const coverImageUrl = `http://localhost:8080/covers/${fileName}`;
    // Add cover image data to the first element of the array
    const updatedCoverImages = [coverImageUrl, ...formData.cover_image];

    // Add data to the database
    const formDataDB = {
      ...formData,
      cover_image: updatedCoverImages, // Update cover_image data
      publisher_id: parseInt(formData.publisher_id),
      type_id: parseInt(formData.type_id),
      status_id: parseInt(formData.status_id)
    };
    try {
      const addDataResponse = await axios.post('http://localhost:8080/add-data', formDataDB);
      console.log(addDataResponse.data);
    } catch (error) {
      console.error('Error adding data:', error);
    }
  };

  const handleChange = (e) => {
    const value = e.target.type === 'number' ? parseFloat(e.target.value) : e.target.value;
    setFormData({ ...formData, [e.target.name]: value });
  };

  const fetchDataFromMyAnimeList = async () => {
    const malId = formData.mal_id;
    if (!malId) {
      console.error('mal_id is empty');
      return;
    }

    try {
      const response = await axios.get(`http://localhost:8080/fetch-mal-data/${malId}`);
      const animeData = response.data;
      console.log(animeData);

      let typeId = '';
      switch (animeData.media_type) {
        case 'manga':
          typeId = '1';
          break;
        case 'light_novel':
          typeId = '2';
          break;
        case 'novel':
          typeId = '3';
          break;
        default:
          typeId = '';
          break;
      }

      // Update state with retrieved data
      setFormData({
        ...formData,
        score: animeData.mean,
        popularity: animeData.num_list_users,
        type_id: typeId
      });
    } catch (error) {
      console.error('Error fetching data from MyAnimeList:', error);
    }
  };

  return (
    <div>
      <h1>Add new catalogue entry</h1>
    <form>
      <input type="file" onChange={handleFileChange} />

      <input type="text" name="title" value={formData.title} onChange={handleChange} placeholder="Title" />

      <select name="type_id" value={formData.type_id} onChange={handleChange}>
        <option value="">Select Type</option>
        {types.map(type => (
          <option key={type.id} value={type.id}>{type.name}</option>
        ))}
      </select>

      <select name="publisher_id" value={formData.publisher_id} onChange={handleChange}>
        <option value="">Select Publisher</option>
        {publishers.map(publisher => (
          <option key={publisher.id} value={publisher.id}>{publisher.name}</option>
        ))}
      </select>

      <select name="status_id" value={formData.status_id} onChange={handleChange}>
        <option value="">Select status</option>
        {status.map(status => (
          <option key={status.id} value={status.id}>{status.name}</option>
        ))}
      </select>

      <input type="number" name="score" value={formData.score} onChange={handleChange} placeholder="Score" />
      <input type="number" name="mal_id" value={formData.mal_id} onChange={handleChange} placeholder="malId" />
      <input type="number" name="popularity" value={formData.popularity} onChange={handleChange} placeholder="popularity" />

      <button type="button" onClick={handleUploadAndAddData}>Upload and Add Data</button>
      <button type="button" onClick={fetchDataFromMyAnimeList}>Fetch Data from MyAnimeList</button>
    </form>
    </div>
  );
}

export default AddTitleForm;
