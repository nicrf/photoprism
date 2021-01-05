package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/photoprism/photoprism/internal/acl"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/internal/i18n"
	"github.com/photoprism/photoprism/internal/query"
	"github.com/photoprism/photoprism/internal/service"
	"github.com/photoprism/photoprism/pkg/txt"
	"net/http"
	"strconv"
)

// GET /api/v1/peoples
func GetPeoples(router *gin.RouterGroup) {
	router.GET("/peoples", func(c *gin.Context) {
		s := Auth(SessionID(c), acl.ResourcePeople, acl.ActionSearch)

		if s.Invalid() {
			AbortUnauthorized(c)
			return
		}

		var f form.PeopleSearch

		err := c.MustBindWith(&f, binding.Form)

		if err != nil {
			AbortBadRequest(c)
			return
		}

		// Guest permissions are limited to shared albums.
		if s.Guest() {
			f.UID = s.Shares.Join(query.Or)
		}

		result, err := query.GetPeoples(0, 50)

		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": txt.UcFirst(err.Error())})
			return
		}

		c.Header("X-Count", strconv.Itoa(len(result)))
		c.Header("X-Limit", strconv.Itoa(f.Count))
		c.Header("X-Offset", strconv.Itoa(f.Offset))

		c.JSON(http.StatusOK, result)
	})
}

// GET /api/v1/peoples/:uid
func GetPeople(router *gin.RouterGroup) {
	router.GET("/peoples/:uid", func(c *gin.Context) {
		s := Auth(SessionID(c), acl.ResourcePeople, acl.ActionRead)

		if s.Invalid() {
			AbortUnauthorized(c)
			return
		}

		id := c.Param("uid")
		a, err := query.PeopleByUID(id)

		if err != nil {
			Abort(c, http.StatusNotFound, i18n.ErrPoepleNotFound)
			return
		}

		c.JSON(http.StatusOK, a)
	})
}

// POST /api/v1/peoples
func CreatePeople(router *gin.RouterGroup) {
	router.POST("/peoples", func(c *gin.Context) {
		s := Auth(SessionID(c), acl.ResourcePeople, acl.ActionCreate)

		if s.Invalid() {
			AbortUnauthorized(c)
			return
		}

		var f form.People

		if err := c.BindJSON(&f); err != nil {
			AbortBadRequest(c)
			return
		}

		a := entity.NewPeople(f.PeopleFullName, f.PeopleUserId, f.PeopleBoD, f.PeopleDeadDate)

		log.Debugf("people: creating %+v %+v", f, a)

		if res := entity.Db().Create(a); res.Error != nil {
			AbortAlreadyExists(c, txt.Quote(a.PeopleFullName))
			return
		}

		event.SuccessMsg(i18n.MsgPeopleCreated)

		UpdateClientConfig()

		PublishPeopleEvent(EntityCreated, a.PeopleUID, c)

		//TODO SaveAlbumAsYaml(*a)

		c.JSON(http.StatusOK, a)
	})
}

// PUT /api/v1/peoples/:uid
func UpdatePeople(router *gin.RouterGroup) {
	router.PUT("/peoples/:uid", func(c *gin.Context) {
		s := Auth(SessionID(c), acl.ResourcePeople, acl.ActionUpdate)

		if s.Invalid() {
			AbortUnauthorized(c)
			return
		}

		uid := c.Param("uid")
		a, err := query.PeopleByUID(uid)

		if err != nil {
			Abort(c, http.StatusNotFound, i18n.ErrAlbumNotFound)
			return
		}

		f, err := form.NewPeople(a)

		if err != nil {
			log.Error(err)
			AbortSaveFailed(c)
			return
		}

		if err := c.BindJSON(&f); err != nil {
			log.Error(err)
			AbortBadRequest(c)
			return
		}

		if err := a.SaveForm(f); err != nil {
			log.Error(err)
			AbortSaveFailed(c)
			return
		}

		UpdateClientConfig()

		event.SuccessMsg(i18n.MsgPeopleSaved)

		PublishPeopleEvent(EntityUpdated, uid, c)

		///TODO SaveAlbumAsYaml(a)

		c.JSON(http.StatusOK, a)
	})
}

// DELETE /api/v1/peoples/:uid
func DeletePeople(router *gin.RouterGroup) {
	router.DELETE("/peoples/:uid", func(c *gin.Context) {
		s := Auth(SessionID(c), acl.ResourcePeople, acl.ActionDelete)

		if s.Invalid() {
			AbortUnauthorized(c)
			return
		}

		conf := service.Config()
		id := c.Param("uid")

		a, err := query.PeopleByUID(id)

		if err != nil {
			Abort(c, http.StatusNotFound, i18n.ErrAlbumNotFound)
			return
		}

		PublishPeopleEvent(EntityDeleted, id, c)

		conf.Db().Delete(&a)

		UpdateClientConfig()

		//TODO SaveAlbumAsYaml(a)

		event.SuccessMsg(i18n.MsgPeopleDeleted, txt.Quote(a.PeopleFullName))

		c.JSON(http.StatusOK, a)
	})
}
