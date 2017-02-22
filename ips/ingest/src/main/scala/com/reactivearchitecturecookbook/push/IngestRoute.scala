package com.reactivearchitecturecookbook.push

import akka.http.scaladsl.model.HttpRequest
import akka.http.scaladsl.server._
import akka.http.scaladsl.unmarshalling.GenericUnmarshallers
import com.nimbusds.jwt.JWTClaimsSet
import com.reactivearchitecturecookbook.ingest.v1m0.IngestedImage

import scala.concurrent.{ExecutionContext, Future}
import scala.util.Try

trait IngestRoute extends Directives with GenericUnmarshallers {

  def authorizeAndExtract(token: String): Try[JWTClaimsSet]

  def extractIngestImage(request: HttpRequest)(implicit ec: ExecutionContext): Future[IngestedImage]

  def publishIngestedImage(claimsSet: JWTClaimsSet, token: String, transactionId: String)(ingestedImage: IngestedImage)(implicit ec: ExecutionContext): Future[Unit]

  private def authorizedRouteHandler(requestContext: RequestContext, token: String, claimsSet: JWTClaimsSet, transactionId: String)(implicit ec: ExecutionContext): Future[RouteResult] = {
    extractIngestImage(requestContext.request)
      .flatMap(publishIngestedImage(claimsSet, token, transactionId))
      .flatMap { _ ⇒ requestContext.complete("") }
      .recoverWith { case x ⇒ requestContext.fail(x) }
  }

  def ingestRoute(implicit ec: ExecutionContext): Route =
    path(RemainingPath) { transactionId ⇒
      post {
        extractCredentials {
          case Some(credentials) if credentials.scheme() == "Bearer" ⇒
            request: RequestContext ⇒
            authorizeAndExtract(credentials.token())
              .map(cs ⇒ authorizedRouteHandler(request, credentials.token(), cs, transactionId.toString()))
              .getOrElse(request.reject(AuthorizationFailedRejection))
          case _ ⇒ _.reject(AuthorizationFailedRejection)
        }
      }
    }

}
