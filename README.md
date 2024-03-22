# genql
Prisma & GraphQL Generator for the Node.js Ecosystem

Inspired by generators from <a href="https://hexdocs.pm/phoenix/Mix.Tasks.Phx.Gen.html" target="_blank">Elixir's Phoenix Framework</a> and <a href="https://guides.rubyonrails.org/command_line.html#bin-rails-generate" target="_blank">Ruby on Rails</a>
# Install
Genql is installable via homebrew if you’re on Mac or Linux. To install, run the following two commands:

```
$ brew install tk04/tap/genql
```

Before using genql, you must have Prisma set up on your project. After that, you can create a model by running:

```
$ genql model [model name] field:type:default_value/"unique" ....
```
Any type value that starts with an upper case character will escape type checking and be inserted in the prisma.schema file as is.

# Example Usage
You can create as many fields as you want, and even generate one-to-one, one-to-many, and many-to-many relationships between models. The third value in the colon-separated argument is reserved for default values. You can also include the “unique” keyword, which will add the “@unique” attribute to a given field.

For example, if we wanted to create a model named “User”, with an auto-incrementing id,  name of type string, and an age of type int, you’d run the following command:
```
$ genql model User id:id:ai name:string age:int
```
This will add the following code block to your prisma.schema file:
```prisma
model User {
	id Int @id @default(autoincrement())
	name String 
	age Int 
}
```
In `id:id:ai`, “id” is the field name, “id” is the type, and “ai” (stands for auto increment) dictates the type of the id field and the default value. The other viable entry for an id type is `id:id:uuid`, which generates an Id of type “String”, with a default value of a UUID string. 

Genql also supports optional values. For example, if you wanted to make the name value from the previously created model optional, you’d do:
```
$ genql model User id:id:ai name:"string?" age:int
```
This will generate the following code: 
```prisma
model User {
	id Int @id @default(autoincrement())
	name String? 
	age Int 
}
```

Let’s also create a Friend model, with a default UUID id field, a unique string email field, a name field of type string, and establish a one-to-one relationship between the User and the Friend models. Here is the command to do such a task: 
```
$ genql model Friend id:id:uuid email:string:unique name:string --OneToOne user:User
```
This command will add the following code to your prisma.schema file:

```prisma
model User {
	id Int @id	@default(autoincrement())
	name String? 
	age Int 
	friend Friend? 
}
model Friend {
	id String @id	@default(uuid())
	email String @unique
	name String 
	user User @relation(fields: [userId], references: [id])
	userId Int @unique
}
```

Notice that the command also adjusted the User model to establish the one-to-one relationship between the models.

# Create GraphQL Resolvers
Other than the model command, genql also supports a “resolvers” command, which will create GraphQL resolvers for a given Prisma model. This command makes a lot of assumptions about your current code organization and structure, so it’s not for everyone. It assumes that you use Type-GraphQL, and organize your resolvers under appname/src/resolvers/. The command also creates a types.ts file with each resolver that includes input & output types to be used in the queries and mutations.

To showcase using it on the Friend model we created previously, run:
```
$ genql resolvers Friend
```
This will create 2 files, both under appname/src/resolvers/Friend. The Index.ts file will contain the queries & mutations, while the types.ts file will contain the types. The “resolvers” command is supposed to be used as a basic entry point to get you started with developing your resolvers. So you’re expected to add error handling and adjust the generated code depending on your needs.

This is what the generated Index.ts file will look like:

```typescript
import { Arg, Ctx, Mutation, Query, Resolver } from "type-graphql";
import { context } from "../context";
import { Friend, createFriendInput, updateFriendInput } from "./types";

@Resolver()
export class FriendResolver {
  @Mutation(() => Friend, { nullable: true })
  deleteFriend(@Ctx() { prisma }: context, @Arg("id") id: string) {
    return prisma.friend.delete({
      where: {
        id: id,
      },
    });
  }
  @Query(() => Friend, { nullable: true })
  getFriend(@Ctx() { prisma }: context, @Arg("id") id: string) {
    return prisma.friend.findFirst({
      where: {
        id: id,
      },
    });
  }
  @Mutation(() => Friend)
  createFriend(
    @Ctx() { prisma }: context,
    @Arg("input") input: createFriendInput
  ) {
    return prisma.friend.create({
      data: {
        ...input,
      },
    });
  }
  @Mutation(() => Friend)
  updateFriend(
    @Ctx() { prisma }: context,
    @Arg("input") input: updateFriendInput
  ) {
    return prisma.friend.update({
      where: { id: input.id },
      data: {
        ...input,
      },
    });
  }
}
```

And this is what the generated types.ts file will look like: 
```typescript
import { Field, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class Friend {
  @Field({ nullable: true })
  id?: string;
  @Field()
  email: string;
  @Field()
  name: string;
  @Field()
  userId: number;
}
@InputType()
export class createFriendInput {
  @Field({ nullable: true })
  id?: string;
  @Field()
  email: string;
  @Field()
  name: string;
  @Field()
  userId: number;
}
@InputType()
export class updateFriendInput {
  @Field()
  id: string;
  @Field({ nullable: true })
  email?: string;
  @Field({ nullable: true })
  name?: string;
  @Field({ nullable: true })
  userId?: number;
}
```
