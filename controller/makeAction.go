package controller

import "github.com/gofiber/fiber/v2"

type Action struct { name string }

func makeAction(name string) *Action {
  return &Action{name: name}
}

func (a *Action) success(code int, msg string, data ...any) *fiber.Map {
  return &fiber.Map{
    "action": a.name,
    "success": true,
    "status_code": code,
    "message": msg,
    "data": data,
  }
}

func (a *Action) fail(code int, msg string) *fiber.Map {
  return &fiber.Map{
    "action": a.name,
    "success": false,
    "status_code": code,
    "message": msg,
  }
}
