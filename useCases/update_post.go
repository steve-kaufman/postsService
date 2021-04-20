package useCases

import (
	"github.com/steve-kaufman/postsService/entities"
	"github.com/steve-kaufman/postsService/interfaces"
)

func UpdatePost(getter interfaces.PostGetter, updater interfaces.PostUpdater, id int, updateData entities.Post) (entities.Post, error) {
	post, err := getter.GetPost(id)
	if err != nil {
		return entities.Post{}, determineError(err)
	}
	return verifyFieldsAndUpdatePost(post, updater, id, updateData)
}

func verifyFieldsAndUpdatePost(original entities.Post, updater interfaces.PostUpdater, id int, updateData entities.Post) (entities.Post, error) {
	err := verifyFields(updateData)
	if err != nil {
		return entities.Post{}, err
	}
	post := updateFields(original, updateData)
	return attemptUpdatePost(updater, post, id)
}

func verifyFields(updateData entities.Post) error {
	if updateData.Likes != 0 || updateData.Dislikes != 0 {
		return ErrCantChangeLikes
	}
	return nil
}

func updateFields(original entities.Post, updateData entities.Post) entities.Post {
	if updateData.Title != "" {
		original.Title = updateData.Title
	}
	if updateData.Content != "" {
		original.Content = updateData.Content
	}
	return original
}

func attemptUpdatePost(updater interfaces.PostUpdater, post entities.Post, id int) (entities.Post, error) {
	err := updater.UpdatePost(id, post)
	if err != nil {
		return entities.Post{}, determineError(err)
	}
	return post, nil
}
